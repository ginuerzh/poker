package poker

import (
	"errors"
	"log"
	"strconv"
	"time"
)

type Occupant struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`

	Pos    int      `json:"index"`
	Chips  int      `json:"chips"`
	Bet    int      `json:"bet"`
	Action string   `json:"action"`
	Cards  []string `json:"cards"`

	conn *Conn
	room *Room

	recv    chan *Message
	Actions chan *Message `json:"-"`
}

func NewOccupant(conn *Conn) *Occupant {
	o := &Occupant{
		conn:    conn,
		recv:    make(chan *Message, 8),
		Actions: make(chan *Message, 1),
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

func (o *Occupant) Room() *Room {
	return o.room
}

func (o *Occupant) Broadcast(message *Message) {
	if o.room == nil {
		return
	}

	for _, oc := range o.room.Occupants {
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
	if timeout < 0 {
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

func (o *Occupant) GetAction(timeout time.Duration) (*Message, error) {
	timer := time.NewTimer(timeout)
	select {
	case m := <-o.Actions:
		return m, nil
	case <-timer.C:
		return nil, errors.New("timeout")
	}
}

func (o *Occupant) Join(rid string) (room *Room) {
	room = GetRoom(rid)
	if room == nil {
		return
	}

	room.lock.Lock()
	defer room.lock.Unlock()

	for pos, _ := range room.Occupants {
		if room.Occupants[pos] == nil {
			room.Occupants[pos] = o
			room.N++
			o.room = room
			o.Pos = pos
			break
		}
	}

	o.Broadcast(&Message{
		Id:     room.Id,
		Type:   MsgPresence,
		From:   o.Id,
		Action: ActJoin,
		Class:  strconv.Itoa(o.Pos),
	})

	o.SendMessage(&Message{
		Id:     room.Id,
		Type:   MsgPresence,
		From:   room.Id,
		Action: ActJoin,
		Class:  strconv.Itoa(o.Pos),
		State:  room,
	})

	return
}

func (o *Occupant) Leave() (room *Room) {
	room = o.room
	if room == nil {
		return
	}

	room.lock.Lock()
	defer room.lock.Unlock()

	o.Broadcast(&Message{
		Id:     room.Id,
		Type:   MsgPresence,
		From:   o.Id,
		Action: ActLeave,
	})

	room.N--
	room.Occupants[o.Pos] = nil
	o.room = nil
	o.Pos = 0

	return
}

func (o *Occupant) Next() *Occupant {
	if o.room == nil {
		return nil
	}

	for i := (o.Pos + 1) % o.room.Cap(); i != o.Pos; i = (i + 1) % o.room.Cap() {
		if o.room.Occupants[i] != nil {
			return o.room.Occupants[i]
		}
	}

	return nil
}
