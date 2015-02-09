package poker

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
)

const (
	MsgIQ       = "iq"
	MsgPresence = "presence"
	MsgMessage  = "message"

	ActPreflop  = "preflop"
	ActFlop     = "flop"
	ActTurn     = "turn"
	ActRiver    = "river"
	ActShowdown = "showdown"

	ActActive = "active"
	ActJoin   = "join"
	ActLeave  = "gone"
	ActBet    = "bet"
	ActButton = "button"
	ActState  = "state"

	ActCall  = "call"
	ActCheck = "check"
	ActRaise = "raise"
	ActFold  = "fold"
	ActAllin = "allin"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

type Message struct {
	Id     string        `json:"id,omitempty"`
	Type   string        `json:"type"`
	From   string        `json:"from,omitempty"`
	To     string        `json:"to,omitempty"`
	Action string        `json:"action"`
	Class  string        `json:"class,omitempty"`
	Body   []interface{} `json:"body,omitempty"`
	State  *Room         `json:"state,omitempty"`
}

type Version struct {
	Id  string `json:"id"`
	Ver string `json:"version"`
}

type Auth struct {
	Mechanism string `json:"mechanism"`
	Text      string `json:"text"`
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
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
		return
	}
	defer ws.Close()

	ws.SetReadLimit(maxMessageSize)
	ws.SetPongHandler(
		func(string) error {
			ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
	conn := NewConn(ws, 256)
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
	if err := conn.WriteJSON(&Error{Code: 0, Err: "success"}); err != nil {
		return
	}

	o := NewOccupant(conn)
	o.Id = strconv.FormatInt(time.Now().UnixNano(), 10)
	o.Name = auth.Text
	o.Chips = 1000

	for {
		message, _ := o.GetMessage(-1)
		if message == nil {
			break
		}

		switch message.Type {
		case MsgIQ:
		case MsgPresence:
			go handlePresence(o, message)
		case MsgMessage:
		}
	}

	o.Leave()
}

func handlePresence(o *Occupant, message *Message) {
	switch message.Action {
	case ActActive:
	case ActJoin:
		room := o.Join("")
		if room == nil {
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
