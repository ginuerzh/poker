package main

import (
	"github.com/ginuerzh/poker"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"
)

func main() {
	c, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	u, err := url.Parse("http://localhost:8000/")
	if err != nil {
		log.Fatal(err)
	}

	ws, _, err := websocket.NewClient(c, u, nil, 1024, 1024)
	if err != nil {
		log.Fatal(err)
	}
	conn := poker.NewConn(ws, 10)

	ver := &poker.Version{
		Ver: "1.0",
	}

	if err := conn.WriteJSON(ver); err != nil {
		log.Fatal(err)
	}

	if err := conn.ReadJSON(ver); err != nil {
		log.Fatal(err)
	}

	auth := poker.Auth{Mechanism: "plain", Text: strconv.FormatInt(time.Now().Unix(), 10)}
	if err := conn.WriteJSON(auth); err != nil {
		log.Fatal(err)
	}

	resp := &poker.Error{}
	if err := conn.ReadJSON(resp); err != nil {
		log.Fatal(err)
	}

	if resp.Code > 0 {
		return
	}

	occupant := poker.NewOccupant(conn)

	msg := &poker.Message{
		Type:   poker.MsgPresence,
		From:   auth.Text,
		Action: poker.ActJoin,
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Fatal(err)
	}

	if err := conn.ReadJSON(msg); err != nil {
		log.Fatal(err)
	}

	log.Println("pos", msg.Class)

	for {

	}
}
