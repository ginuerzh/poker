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
	Chips     []int       `json:"chips"`
	Bet       int         `json:"bet"`
	N         int         `json:"n"`
	remain    int
	EndChan   chan int `json:"-"`
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
		Chips:     make([]int, n),
		SB:        sb,
		BB:        bb,
		Pot:       make([]int, 1),
		Timeout:   30,
		lock:      sync.Mutex{},
		deck:      NewDeck(),
		EndChan:   make(chan int),
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
	/*
		for pos, _ := range room.Occupants {
			p := room.Occupants[pos]
			if p != nil && p.Id == o.Id {
				o.Bet = p.Bet
				o.Cards = p.Cards
				o.Hand = p.Hand
				o.Action = p.Action
				// o.gone = p.gone
				o.Pos = p.Pos
				o.Room = p.Room
				o.Chips = p.Chips

				room.Occupants[pos] = o

				return o.Pos
			}
		}
	*/
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
	if len(o.Cards) > 0 {
		room.remain--
	}
	if room.remain <= 1 {
		select {
		case room.EndChan <- 0:
		default:
		}
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
		if room.Occupants[i] != nil && !f(room.Occupants[i]) {
			return
		}
	}

	// end
	if room.Occupants[i] != nil {
		f(room.Occupants[i])
	}
}

func (room *Room) start() {
	var dealer *Occupant

	room.Each(0, func(o *Occupant) bool {
		if o.Chips < room.BB {
			o.Leave()
		}
		return true
	})

	// Select Dealer
	button := room.Button - 1
	room.Each((button+1)%room.Cap(), func(o *Occupant) bool {
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

	room.deck.Shuffle()

	// Small Blind
	sb := dealer.Next()
	if room.N == 2 { // one-to-one
		sb = dealer
	}
	// Big Blind
	bb := sb.Next()
	bbPos := bb.Pos

	room.Pot = nil
	for i, _ := range room.Chips {
		room.Chips[i] = 0
	}
	room.Bet = 0
	room.Cards = nil
	room.remain = 0
	room.Each(0, func(o *Occupant) bool {
		o.Bet = 0
		o.Cards = []Card{room.deck.Take(), room.deck.Take()}
		o.Hand = 0
		//o.Action = ActReady
		o.Action = ""
		room.remain++

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
	room.calc()

	// Round 2 : Flop
	room.ready()
	room.Cards = []Card{
		room.deck.Take(),
		room.deck.Take(),
		room.deck.Take(),
	}
	room.Each(0, func(o *Occupant) bool {
		var hand [5]Card
		if len(o.Cards) > 0 {
			cards := hand[0:0]
			cards = append(cards, o.Cards...)
			cards = append(cards, room.Cards...)
			o.Hand = Eva5Hand(hand)
		}
		o.SendMessage(&Message{
			From:   room.Id,
			Type:   MsgPresence,
			Action: ActFlop,
			Class:  fmt.Sprintf("%s,%s,%s,%d", room.Cards[0], room.Cards[1], room.Cards[2], o.Hand>>16),
		})

		return true
	})
	room.action(0)
	if room.remain <= 1 {
		goto showdown
	}
	room.calc()

	// Round 3 : Turn
	room.ready()
	room.Cards = append(room.Cards, room.deck.Take())
	room.Each(0, func(o *Occupant) bool {
		var hand [6]Card
		if len(o.Cards) > 0 {
			cards := hand[0:0]
			cards = append(cards, o.Cards...)
			cards = append(cards, room.Cards...)
			o.Hand = Eva6Hand(hand)
		}
		o.SendMessage(&Message{
			From:   room.Id,
			Type:   MsgPresence,
			Action: ActTurn,
			Class:  fmt.Sprintf("%s,%d", room.Cards[3], o.Hand>>16),
		})

		return true
	})
	room.action(0)
	if room.remain <= 1 {
		goto showdown
	}
	room.calc()

	// Round 4 : River
	room.ready()
	room.Cards = append(room.Cards, room.deck.Take())
	room.Each(0, func(o *Occupant) bool {
		var hand [7]Card
		if len(o.Cards) > 0 {
			cards := hand[0:0]
			cards = append(cards, o.Cards...)
			cards = append(cards, room.Cards...)
			o.Hand = Eva7Hand(hand)
		}
		o.SendMessage(&Message{
			From:   room.Id,
			Type:   MsgPresence,
			Action: ActRiver,
			Class:  fmt.Sprintf("%s,%d", room.Cards[4], o.Hand>>16),
		})

		return true
	})
	room.action(0)

showdown:
	room.showdown()
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
	if pos == 0 { // start from left hand of button
		pos = (room.Button)%room.Cap() + 1
	}

	for {
		raised := 0

		room.Each(pos-1, func(o *Occupant) bool {
			if room.remain <= 1 {
				return false
			}

			if o.Pos == skip || o.Chips == 0 || len(o.Cards) == 0 {
				return true
			}

			room.Broadcast(&Message{
				From:   room.Id,
				Type:   MsgPresence,
				Action: ActAction,
				Class:  fmt.Sprintf("%d,%d", o.Pos, room.Bet),
			})

			msg, _ := o.GetAction(time.Duration(room.Timeout) * time.Second)
			if room.remain <= 1 {
				return false
			}

			n := 0
			// timeout or leave
			if msg == nil || len(msg.Class) == 0 {
				n = -1
			} else {
				n, _ = strconv.Atoi(msg.Class)
			}

			if room.betting(o.Pos, n) {
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
}

func (room *Room) calc() (pots []handPot) {
	pots = calcPot(room.Chips)
	room.Pot = nil
	var ps []string
	for _, pot := range pots {
		room.Pot = append(room.Pot, pot.Pot)
		ps = append(ps, strconv.Itoa(pot.Pot))
	}

	room.Broadcast(&Message{
		From:   room.Id,
		Type:   MsgPresence,
		Action: ActPot,
		Class:  strings.Join(ps, ","),
	})

	return
}

func (room *Room) showdown() {
	pots := room.calc()

	for i, _ := range room.Chips {
		room.Chips[i] = 0
	}

	room.lock.Lock()
	defer room.lock.Unlock()

	for _, pot := range pots {
		maxHand := 0
		for _, pos := range pot.OPos {
			o := room.Occupants[pos-1]
			if o != nil && o.Hand > maxHand {
				maxHand = o.Hand
			}
		}

		var winners []int

		for _, pos := range pot.OPos {
			o := room.Occupants[pos-1]
			if o != nil && o.Hand == maxHand && len(o.Cards) > 0 {
				winners = append(winners, pos)
			}
		}

		if len(winners) == 0 {
			fmt.Println("!!!no winners!!!")
			return
		}

		for _, winner := range winners {
			room.Chips[winner-1] += pot.Pot / len(winners)
		}
		room.Chips[winners[0]-1] += pot.Pot % len(winners) // odd chips
	}

	for i, _ := range room.Chips {
		if room.Occupants[i] != nil {
			room.Occupants[i].Chips += room.Chips[i]
		}
	}
}

func (room *Room) ready() {
	room.Bet = 0
	room.lock.Lock()
	defer room.lock.Unlock()

	room.Each(0, func(o *Occupant) bool {
		o.Bet = 0
		/*
			if o.Action == ActAllin || o.Action == ActFold || o.Action == "" {
				return true
			}
			o.Action = ActReady
		*/
		return true
	})

}

func (room *Room) betting(pos, n int) (raised bool) {
	if pos <= 0 {
		return
	}

	room.lock.Lock()
	defer room.lock.Unlock()

	o := room.Occupants[pos-1]
	if o == nil {
		return
	}
	raised = o.Betting(n)
	if o.Action == ActFold {
		room.remain--
	}

	room.Broadcast(&Message{
		Id:     room.Id,
		Type:   MsgPresence,
		From:   o.Id,
		Action: ActBet,
		Class:  o.Action + "," + strconv.Itoa(o.Bet) + "," + strconv.Itoa(o.Chips),
	})

	return
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
