package main

import (
	"bufio"
	"fmt"
	"github.com/ginuerzh/poker"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	readWait = 10 * time.Second
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

	o := poker.NewOccupant(conn)

	go cmdLoop(o)

	for {
		message, _ := o.GetMessage(-1)
		if message == nil {
			break
		}

		switch message.Type {
		case poker.MsgPresence:
			handlePresence(o, message)
		}
	}

}

func handlePresence(o *poker.Occupant, message *poker.Message) {
	switch message.Action {
	case poker.ActButton:

	}
}

func cmdLoop(o *poker.Occupant) {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("poker> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.ToLower(strings.Trim(cmd, " \n"))

		if len(cmd) == 0 {
			continue
		}
		switch cmd[0] {
		case 'j':
			o.SendMessage(&poker.Message{
				Type:   poker.MsgPresence,
				Action: poker.ActJoin,
			})

		case 'l':
			o.SendMessage(&poker.Message{
				Type:   poker.MsgPresence,
				Action: poker.ActLeave,
			})
		case 'q':
			return
		default:
			bet, _ := strconv.ParseInt(cmd, 10, 32)
			o.SendMessage(&poker.Message{
				Type:   poker.MsgPresence,
				Action: poker.ActBet,
				Class:  strconv.FormatInt(bet, 10),
			})
		}
	}
}
