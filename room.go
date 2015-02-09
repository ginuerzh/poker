package poker

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

type Room struct {
	Id        string      `json:"id"`
	SB        int         `json:"sb"`
	BB        int         `json:"bb"`
	Cards     []string    `json:"cards"`
	Mainpot   int         `json:"mainpot"`
	Sidepot   []int       `json"sidepot"`
	Button    int         `json:"button"`
	Occupants []*Occupant `json:"occupants"`
	Bet       int         `json:"bet"`
	N         int         `json:"n"`
	lock      *sync.Mutex
}

func NewRoom(id string, n int, sb, bb int) *Room {
	if n <= 0 {
		n = 10 // default 10 occupants
	}

	room := &Room{
		Id:        id,
		Occupants: make([]*Occupant, n),
		SB:        sb,
		BB:        bb,
		lock:      &sync.Mutex{},
	}
	go func() {
		for {
			room.start()
			time.Sleep(1 * time.Second)
		}
	}()

	return room
}

func (room *Room) Cap() int {
	return len(room.Occupants)
}

func (room *Room) Len() int {
	return room.N
}

func (room *Room) Broadcast(message *Message) {
	for _, o := range room.Occupants {
		if o != nil {
			o.SendMessage(message)
		}
	}
}

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

func (room *Room) action(start int) {
	skip := 0

	for {
		raised := 0

		room.Each(start, func(o *Occupant) bool {
			if o == nil || o.Action == ActFold || o.Chips == 0 || o.Pos+1 == skip {
				return true
			}

			msg, _ := o.GetAction(readWait)

			m := &Message{
				Id:     room.Id,
				Type:   MsgPresence,
				From:   o.Id,
				Action: ActBet,
				Class:  "-1",
			}

			if msg == nil || len(msg.Class) == 0 {
				o.Action = ActFold
				room.Broadcast(m)
				return true
			}

			m.Class = msg.Class
			room.Broadcast(m)

			n, _ := strconv.Atoi(msg.Class)
			if n < 0 {
				o.Action = ActFold
			} else if n == 0 {
				o.Action = ActCheck
			} else if n <= room.Bet {
				o.Action = ActCall
				o.Chips -= (n - o.Bet)
				o.Bet = n
			} else {
				o.Action = ActRaise
				o.Chips -= (n - o.Bet)
				room.Bet = n
				o.Bet = room.Bet
				raised = o.Pos + 1 // bet raised
				return false
			}

			return true
		})

		if raised == 0 {
			break
		}

		start = raised - 1
		skip = start
	}

	room.Each(0, func(o *Occupant) bool {
		if o != nil {
			room.Mainpot += o.Bet
			o.Bet = 0
			o.Action = ""
		}
		return true
	})
}

func (room *Room) start() {

	var dealer *Occupant
	// Select Dealer
	room.Each((room.Button+1)%room.Cap(), func(o *Occupant) bool {
		if o == nil {
			return true
		}
		room.Button = o.Pos
		dealer = o
		return false
	})

	if dealer == nil || room.Len() < 2 {
		return
	}

	room.Broadcast(&Message{
		Id:     room.Id,
		Type:   MsgPresence,
		Action: ActButton,
		Class:  strconv.Itoa(room.Button),
	})

	// Small Blind
	sb := dealer.Next()
	if room.Len() == 2 {
		sb = dealer
	}
	sb.Action = ActCall
	sb.Bet = room.SB
	sb.Chips -= sb.Bet

	// Big Blind
	bb := sb.Next()
	bb.Action = ActRaise
	bb.Bet = room.BB
	bb.Chips -= bb.Bet

	room.Bet = room.BB

	// Round 1 : preflop
	room.Each(sb.Pos, func(o *Occupant) bool {
		if o == nil {
			return true
		}
		card1, card2 := "card1", "card2"
		o.Cards = []string{card1, card2}

		o.SendMessage(&Message{
			Id:     room.Id,
			Type:   MsgPresence,
			Action: ActPreflop,
			Class:  "card1,card2",
		})
		return true
	})
	// Under the Gun
	utg := bb.Next()
	room.action(utg.Pos)

	// Round 2 : Flop
	room.Cards = []string{"card1", "card2", "card3"}
	room.Broadcast(&Message{
		Id:     room.Id,
		Type:   MsgPresence,
		Action: ActFlop,
		Class:  strings.Join(room.Cards, ","),
	})
	room.action(sb.Pos)

	// Round 3 : Turn
	room.Cards = append(room.Cards, "card4")
	room.Broadcast(&Message{
		Id:     room.Id,
		Type:   MsgPresence,
		Action: ActTurn,
		Class:  room.Cards[3],
	})
	room.action(sb.Pos)

	// Round 4 : River
	room.Cards = append(room.Cards, "card5")
	room.Broadcast(&Message{
		Id:     room.Id,
		Type:   MsgPresence,
		Action: ActRiver,
		Class:  room.Cards[4],
	})
	room.action(sb.Pos)

	// Final : Showdown

}

var (
	rooms = map[string]*Room{"": NewRoom("", 10, 5, 10)}
)

func GetRoom(rid string) *Room {
	if _, ok := rooms[rid]; ok {
		return rooms[rid]
	}

	return nil
}
