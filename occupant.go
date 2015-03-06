package poker

import (
	"errors"
	"log"
	//"strconv"
	"time"
)

type Occupant struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Profile string `json:"profile"`
	Level   int    `json:"level"`
	Chips   int    `json:"chips"`

	Pos    int    `json:"index,omitempty"`
	Bet    int    `json:"bet,omitempty"`
	Action string `json:"action,omitempty"`
	Cards  []Card `json:"cards,omitempty"`
	Hand   int    `json:"hand,omitempty"`

	conn *Conn
	Room *Room `json:"-"`

	recv    chan *Message
	Actions chan *Message `json:"-"`
	timer   *time.Timer   // action timer
}

func NewOccupant(id string, conn *Conn) *Occupant {
	o := &Occupant{
		Id:      id,
		conn:    conn,
		recv:    make(chan *Message, 128),
		Actions: make(chan *Message),
		Profile: "http://172.24.222.42:8082/4,0aaaffedfc7b",
	}

	go func() {
		for {
			m := &Message{}
			if err := o.conn.ReadJSON(m); err != nil {
				close(o.recv)
				o.recv = nil
				return
			}
			select {
			case o.recv <- m:
			default:
				log.Println("dropped")
			}
		}
	}()

	return o
}

func (o *Occupant) Broadcast(message *Message) {
	if o.Room == nil {
		return
	}

	for _, oc := range o.Room.Occupants {
		if oc != nil && oc != o {
			oc.SendMessage(message)
		}
	}
}

func (o *Occupant) SendMessage(message *Message) error {
	return o.conn.WriteJSON(message)
}

func (o *Occupant) SendError(code int, err string) error {
	return o.conn.WriteJSON(NewError(code, err))
}

func (o *Occupant) GetMessage(timeout time.Duration) (*Message, error) {
	if o.recv == nil {
		return nil, errors.New("channel closed")
	}
	if timeout <= 0 {
		m := <-o.recv
		return m, nil
	}

	timer := time.NewTimer(timeout)
	select {
	case m := <-o.recv:
		return m, nil
	case <-timer.C:
		return nil, errors.New("timeout")
	}
}

func (o *Occupant) Betting(n int) (raised bool) {
	room := o.Room
	if room == nil {
		return
	}

	if n < 0 {
		o.Action = ActFold
		o.Cards = nil
		o.Hand = 0
		n = 0
	} else if n == 0 {
		o.Action = ActCheck
	} else if n+o.Bet <= room.Bet {
		o.Action = ActCall
		o.Chips -= n
		o.Bet += n
	} else {
		o.Action = ActRaise
		o.Chips -= n
		o.Bet += n
		room.Bet = o.Bet
		raised = true
	}
	if o.Chips == 0 {
		o.Action = ActAllin
	}
	room.Chips[o.Pos-1] += n

	return
}

func (o *Occupant) GetAction(timeout time.Duration) (*Message, error) {
	o.timer = time.NewTimer(timeout)

	select {
	case m := <-o.Actions:
		return m, nil
	case <-o.Room.EndChan:
		return nil, nil
	case <-o.timer.C:
		return nil, errors.New("timeout")
	}
}

func (o *Occupant) Join(rid string) (room *Room) {
	room = GetRoom(rid)
	if room == nil {
		return
	}

	o.Bet = 0
	o.Cards = nil
	o.Hand = 0
	o.Action = ""
	o.Pos = 0
	o.Room = nil

	room.AddOccupant(o)

	o.Broadcast(&Message{
		From:     room.Id,
		Type:     MsgPresence,
		Action:   ActJoin,
		Occupant: o,
	})
	o.SendMessage(&Message{
		From:   room.Id,
		Type:   MsgPresence,
		Action: ActState,
		Room:   room,
	})

	return
}

func (o *Occupant) Leave() (room *Room) {
	room = o.Room
	if room == nil {
		return
	}

	room.DelOccupant(o)

	o.Broadcast(&Message{
		From:     room.Id,
		Type:     MsgPresence,
		Action:   ActLeave,
		Occupant: o,
	})

	o.Bet = 0
	o.Cards = nil
	o.Hand = 0
	o.Action = ""
	o.Pos = 0
	o.Room = nil
	if o.timer != nil {
		o.timer.Reset(0)
	}

	return
}

func (o *Occupant) Next() *Occupant {
	room := o.Room
	if room == nil {
		return nil
	}

	for i := (o.Pos) % room.Cap(); i != o.Pos-1; i = (i + 1) % room.Cap() {
		if room.Occupants[i] != nil {
			return room.Occupants[i]
		}
	}

	return nil
}
