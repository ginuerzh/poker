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

	ActPreflop = "preflop"
	ActFlop    = "flop"
	ActTurn    = "turn"
	ActRiver   = "river"

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
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
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

	ver := &Version{}
	if err := conn.ReadJSON(ver); err != nil {
		return
	}
	if err := conn.WriteJSON(ver); err != nil {
		return
	}

	auth := &Auth{}
	if err := conn.ReadJSON(auth); err != nil {
		return
	}
	if err := conn.WriteJSON(&Error{Code: 0, Err: "success"}); err != nil {
		return
	}

	o := NewOccupant(conn)
	o.Id = strconv.FormatInt(time.Now().UnixNano(), 10)
	o.Name = auth.Text
	o.Chips = 1000

	o.HandleMessage()
}
