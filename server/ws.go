package poker

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	MsgIQ       = "iq"
	MsgPresence = "presence"
	MsgMessage  = "message"

	ActGet    = "get"
	ActSet    = "set"
	ActResult = "result"

	ActPreflop  = "preflop"
	ActFlop     = "flop"
	ActTurn     = "turn"
	ActRiver    = "river"
	ActShowdown = "showdown"
	ActPot      = "pot"

	ActActive = "active"
	ActJoin   = "join"
	ActLeave  = "gone"
	ActBet    = "bet"
	ActButton = "button"
	ActState  = "state"

	ActAction = "action"
	ActReady  = "ready"
	ActCall   = "call"
	ActCheck  = "check"
	ActRaise  = "raise"
	ActFold   = "fold"
	ActAllin  = "allin"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

type Message struct {
	Id       string    `json:"id,omitempty"`
	Type     string    `json:"type"`
	From     string    `json:"from,omitempty"`
	To       string    `json:"to,omitempty"`
	Action   string    `json:"action"`
	Class    string    `json:"class,omitempty"`
	Occupant *Occupant `json:"occupant,omitempty"`
	Room     *Room     `json:"room,omitempty"`
	Rooms    []*Room   `json:"rooms,omitempty"`
}

type Version struct {
	//Id  string `json:"id"`
	Ver string `json:"version"`
}

type Auth struct {
	Mechanism string `json:"mechanism"`
	Text      string `json:"text"`
}

type AuthResp struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
	Chips int    `json:"chips"`
}

type Error struct {
	Code int    `json:"code,omitempty"`
	Err  string `json:"error,omitempty"`
}

func NewError(code int, err string) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

func PokerHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer ws.Close()

	ws.SetReadLimit(maxMessageSize)
	ws.SetPongHandler(
		func(string) error {
			ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
	conn := NewConn(ws, 128)
	defer conn.Close()

	ver := &Version{}
	if err := conn.ReadJSONTimeout(ver, readWait); err != nil {
		return
	}
	if err := conn.WriteJSON(ver); err != nil {
		return
	}

	auth := &Auth{}
	if err := conn.ReadJSONTimeout(auth, readWait); err != nil {
		return
	}

	o := NewOccupant(strconv.FormatInt(time.Now().Unix(), 10), conn)
	o.Name = auth.Text
	o.Chips = 10000

	if err := conn.WriteJSON(o); err != nil {
		return
	}

	for {
		message, _ := o.GetMessage(0)
		if message == nil {
			break
		}

		switch message.Type {
		case MsgIQ:
			go handleIQ(o, message)
		case MsgPresence:
			go handlePresence(o, message)
		case MsgMessage:
		}
	}

	o.Leave()
	log.Println(o.Name, "disconnected.")
}

func handlePresence(o *Occupant, message *Message) {
	switch message.Action {
	case ActJoin:
		if room := o.Join(message.To); room == nil {
			o.SendError(1, "room not found")
			return
		}
	case ActLeave:
		o.Leave()
	case ActBet:
		select {
		case o.Actions <- message:
		default:
		}
	}
}

func handleIQ(o *Occupant, message *Message) {
	switch message.Action {
	case ActSet:
		switch message.Class {
		case "room":
			room := NewRoom("", 9, 5, 10)

			if message.Room != nil {
				if message.Room.SB > 0 {
					room.SB = message.Room.SB
				}
				if message.Room.BB > 0 {
					room.BB = message.Room.BB
				}
				if message.Room.Timeout > 0 {
					room.Timeout = message.Room.Timeout
				}

				if message.Room.Max > 0 && message.Room.Max <= MaxN {
					room.Max = message.Room.Max
					room.Occupants = room.Occupants[:room.Max]
					room.Chips = room.Chips[:room.Max]
				}

				SetRoom(room)
			}
			o.SendMessage(&Message{
				Type:   message.Type,
				Id:     message.Id,
				From:   room.Id,
				Action: ActResult,
				Class:  message.Class,
				Room:   room,
			})
		}

	case ActGet:
		switch message.Class {
		case "room":
			o.SendMessage(&Message{
				Type:   message.Type,
				Id:     message.Id,
				From:   message.To,
				Action: ActResult,
				Class:  message.Class,
				Room:   GetRoom(message.To),
			})
		case "roomlist":
			rooms := []*Room{}
			for _, room := range Rooms() {
				rooms = append(rooms, &Room{
					Id:  room.Id,
					SB:  room.SB,
					BB:  room.BB,
					N:   room.N,
					Max: room.Max,
				})
			}
			o.SendMessage(&Message{
				Type:   message.Type,
				Id:     message.Id,
				Action: ActResult,
				Class:  message.Class,
				Rooms:  rooms,
			})
		}
	}
}
