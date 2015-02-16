package poker

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	actionWait = 20 * time.Second
)

type Room struct {
	Id        string      `json:"id"`
	SB        int         `json:"sb"`
	BB        int         `json:"bb"`
	Cards     []Card      `json:"cards"`
	Pot       []int       `json:"pot"`
	Timeout   int         `json:"timeout"`
	Button    int         `json:"button"`
	Occupants []*Occupant `json:"occupants"`
	Bet       int         `json:"bet"`
	N         int         `json:"n"`
	remain    int
	lock      sync.Mutex
	deck      *Deck
}

func NewRoom(id string, n int, sb, bb int) *Room {
	if n <= 0 {
		n = 9 // default 9 occupants
	}

	room := &Room{
		Id:        id,
		Occupants: make([]*Occupant, n),
		SB:        sb,
		BB:        bb,
		Pot:       make([]int, 1),
		Timeout:   20,
		lock:      sync.Mutex{},
		deck:      NewDeck(),
	}
	go func() {
		for {
			room.start()
			time.Sleep(5 * time.Second)
		}
	}()

	return room
}

func (room *Room) Cap() int {
	return len(room.Occupants)
}

func (room *Room) Occupant(id string) *Occupant {
	for _, o := range room.Occupants {
		if o != nil && o.Id == id {
			return o
		}
	}

	return nil
}

func (room *Room) AddOccupant(o *Occupant) int {
	room.lock.Lock()
	defer room.lock.Unlock()

	for pos, _ := range room.Occupants {
		if room.Occupants[pos] == nil {
			room.Occupants[pos] = o
			room.N++
			o.Room = room
			o.Pos = pos + 1
			break
		}
	}

	return o.Pos
}

func (room *Room) DelOccupant(o *Occupant) {
	if o == nil || o.Pos == 0 {
		return
	}

	room.lock.Lock()
	defer room.lock.Unlock()

	room.Occupants[o.Pos-1] = nil
	room.N--
	if o.Action == ActReady || len(o.Cards) > 0 {
		room.remain--
	}
}

func (room *Room) Broadcast(message *Message) {
	for _, o := range room.Occupants {
		if o != nil {
			o.SendMessage(message)
		}
	}
}

// start starts from 0
func (room *Room) Each(start int, f func(o *Occupant) bool) {
	end := (room.Cap() + start - 1) % room.Cap()
	i := start
	for ; i != end; i = (i + 1) % room.Cap() {
		if !f(room.Occupants[i]) {
			return
		}
	}

	// end
	f(room.Occupants[i])
}

func (room *Room) start() {
	var dealer *Occupant

	room.deck.Shuffle()

	room.Each(0, func(o *Occupant) bool {
		if o != nil && o.Chips < room.BB {
			o.Leave()
		}
		return true
	})

	// Select Dealer
	room.Each((room.Button+1)%room.Cap(), func(o *Occupant) bool {
		if o == nil || o.Chips < room.BB {
			return true
		}
		room.Button = o.Pos
		dealer = o
		return false
	})

	if dealer == nil {
		return
	}

	room.lock.Lock()
	if room.N < 2 {
		room.lock.Unlock()
		return
	}
	// Small Blind
	sb := dealer.Next()
	if room.N == 2 { // one-to-one
		sb = dealer
	}
	// Big Blind
	bb := sb.Next()
	bbPos := bb.Pos

	room.Pot = make([]int, 1)
	room.Bet = 0
	room.Cards = nil
	room.remain = 0
	room.Each(0, func(o *Occupant) bool {
		if o != nil {
			o.Bet = 0
			o.Cards = nil
			o.Hand = 0
			o.Action = ActReady
			room.remain++
		}
		return true
	})
	room.lock.Unlock()

	room.Broadcast(&Message{
		From:   room.Id,
		Type:   MsgPresence,
		Action: ActButton,
		Class:  strconv.Itoa(room.Button),
	})

	room.betting(sb.Pos, room.SB)
	room.betting(bb.Pos, room.BB)

	// Round 1 : preflop
	room.Each(sb.Pos-1, func(o *Occupant) bool {
		if o == nil {
			return true
		}
		o.Cards = []Card{room.deck.Take(), room.deck.Take()}

		o.SendMessage(&Message{
			From:   room.Id,
			Type:   MsgPresence,
			Action: ActPreflop,
			Class:  o.Cards[0].String() + "," + o.Cards[1].String(),
		})
		return true
	})

	room.action(bbPos%room.Cap() + 1)
	if room.remain <= 1 {
		goto showdown
	}

	// Round 2 : Flop
	room.ready()
	room.Cards = []Card{
		room.deck.Take(),
		room.deck.Take(),
		room.deck.Take(),
	}
	room.Broadcast(&Message{
		From:   room.Id,
		Type:   MsgPresence,
		Action: ActFlop,
		Class:  room.Cards[0].String() + "," + room.Cards[1].String() + "," + room.Cards[2].String(),
	})
	room.action(0)
	if room.remain <= 1 {
		goto showdown
	}

	// Round 3 : Turn
	room.ready()
	room.Cards = append(room.Cards, room.deck.Take())
	room.Broadcast(&Message{
		From:   room.Id,
		Type:   MsgPresence,
		Action: ActTurn,
		Class:  room.Cards[3].String(),
	})
	room.action(0)
	if room.remain <= 1 {
		goto showdown
	}

	// Round 4 : River
	room.ready()
	room.Cards = append(room.Cards, room.deck.Take())
	room.Broadcast(&Message{
		From:   room.Id,
		Type:   MsgPresence,
		Action: ActRiver,
		Class:  room.Cards[4].String(),
	})
	room.action(0)

	room.Each(0, func(o *Occupant) bool {
		if o == nil || len(o.Cards) == 0 {
			return true
		}

		var hand [7]Card
		cards := hand[0:0:7]
		cards = append(cards, o.Cards...)
		cards = append(cards, room.Cards...)

		v := Eva7Hand(hand)
		o.Hand = HandRank(v)
		return true
	})

showdown:

	// Final : Showdown
	room.Broadcast(&Message{
		From:   room.Id,
		Type:   MsgPresence,
		Action: ActShowdown,
		Room:   room,
	})
}

func (room *Room) action(pos int) {
	skip := 0
	if pos == 0 { // start from button left
		pos = (room.Button)%room.Cap() + 1
	}

	for {
		raised := 0

		room.Each(pos-1, func(o *Occupant) bool {
			if room.remain <= 1 {
				return false
			}
			if o == nil || o.Pos == skip || o.Chips == 0 || len(o.Cards) == 0 {
				return true
			}

			room.Broadcast(&Message{
				From:   room.Id,
				Type:   MsgPresence,
				Action: ActAction,
				Class:  fmt.Sprintf("%d,%d", o.Pos, room.Bet),
			})
			msg, _ := o.GetAction(time.Duration(room.Timeout) * time.Second)

			n := 0
			// timeout or leave
			if msg == nil || len(msg.Class) == 0 {
				n = -1
			} else {
				n, _ = strconv.Atoi(msg.Class)
			}

			room.betting(o.Pos, n)

			if o.Action == ActRaise {
				raised = o.Pos
				return false
			}

			return true
		})

		if raised == 0 {
			break
		}

		pos = raised
		skip = pos
	}

	var pots []string
	for _, pot := range room.Pot {
		pots = append(pots, strconv.Itoa(pot))
	}
	room.Broadcast(&Message{
		From:   room.Id,
		Type:   MsgPresence,
		Action: ActPot,
		Class:  strings.Join(pots, ","),
	})
}

func (room *Room) ready() {
	room.Bet = 0
	room.lock.Lock()
	defer room.lock.Unlock()

	room.Each(0, func(o *Occupant) bool {
		if o == nil {
			return true
		}
		o.Bet = 0
		o.Action = ActReady
		return true
	})
}

func (room *Room) betting(pos, n int) {
	if pos <= 0 {
		return
	}

	room.lock.Lock()
	o := room.Occupants[pos-1]
	if o == nil {
		return
	}
	o.Betting(n)
	if o.Action == ActFold { // fold
		room.remain--
	}
	room.lock.Unlock()

	room.Broadcast(&Message{
		Id:     room.Id,
		Type:   MsgPresence,
		From:   o.Id,
		Action: ActBet,
		Class:  strconv.Itoa(n),
	})
}

var (
	rooms = map[string]*Room{"0": NewRoom("0", 9, 5, 10)}
)

func GetRoom(rid string) *Room {
	if _, ok := rooms[rid]; ok {
		return rooms[rid]
	}

	return nil
}
