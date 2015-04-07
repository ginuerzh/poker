package poker

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Poker struct {
	WebRoot string
	Addr    string
	OnAuth  func(conn *Conn, mechanism, text string) (*Occupant, error)
	OnExit  func(o *Occupant)
}

func (p *Poker) ListenAndServe() error {
	http.HandleFunc("/ws", p.pokerHandler)
	http.Handle("/", http.FileServer(http.Dir(p.WebRoot)))
	return http.ListenAndServe(p.Addr, nil)
}

func (p *Poker) pokerHandler(w http.ResponseWriter, r *http.Request) {
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

	var o *Occupant
	if p.OnAuth != nil {
		o, err = p.OnAuth(conn, auth.Mechanism, auth.Text)
		if err != nil {
			log.Println(err)
			//conn.WriteJSON(err)
			//return
			o = NewOccupant(strconv.FormatInt(time.Now().Unix(), 10), conn)
			o.Name = auth.Text
		}
		o.Chips = 10000
	}

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

	if p.OnExit != nil {
		p.OnExit(o)
	}
	log.Println(o.Name, "disconnected.")
}
