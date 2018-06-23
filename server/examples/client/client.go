package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ginuerzh/poker/server"
	"github.com/gorilla/websocket"
)

const (
	readWait = 10 * time.Second
	version  = "1.0"
)

func randName() string {
	rand.Seed(time.Now().Unix())
	var b []byte
	for i := 0; i < 5; i++ {
		b = append(b, byte(rand.Intn(26)+97))
	}
	b[0] -= 32
	return string(b)
}

func main() {
	c, err := net.Dial("tcp", "localhost:8989")
	if err != nil {
		log.Fatal(err)
	}
	u, err := url.Parse("ws://localhost:8989/ws")
	if err != nil {
		log.Fatal(err)
	}

	ws, _, err := websocket.NewClient(c, u, nil, 1024, 1024)
	if err != nil {
		log.Fatal(err)
	}
	conn := poker.NewConn(ws, 10)

	ver := &poker.Version{
		Ver: version,
	}

	if err := conn.WriteJSON(ver); err != nil {
		log.Fatal(err)
	}

	if err := conn.ReadJSON(ver); err != nil {
		log.Fatal(err)
	}

	auth := poker.Auth{Mechanism: "plain", Text: randName()}
	if err := conn.WriteJSON(auth); err != nil {
		log.Fatal(err)
	}

	resp := &poker.AuthResp{}
	if err := conn.ReadJSON(resp); err != nil {
		log.Fatal(err)
	}

	o := poker.NewOccupant(resp.Id, conn)
	o.Name = resp.Name
	o.Chips = resp.Chips

	fmt.Printf("%s(%s) %d\n", o.Id, o.Name, o.Chips)

	go handleMessage(o)

	cmdLoop(o)

}

func handleMessage(o *poker.Occupant) {
	for {
		message, _ := o.GetMessage(0)
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
	case poker.ActState:
		o.Room = message.Room
		fmt.Printf("Enter room, %d Occupants\n", o.Room.N)
	case poker.ActJoin:
		occupant := message.Occupant
		o.Room.Occupants[occupant.Pos-1] = occupant
		o.Room.N++
		fmt.Printf("%s(%s) Join.\n", occupant.Id, occupant.Name)
	case poker.ActLeave:
		occupant := message.Occupant
		o.Room.Occupants[occupant.Pos-1] = nil
		o.Room.N--
		if occupant.Id == o.Id {
			o.Room = nil
			o.Pos = 0
			fmt.Println("You are kicked.")
		} else {
			fmt.Printf("%s(%s) Leave.\n", occupant.Id, occupant.Name)
		}
	case poker.ActButton:
		pos, _ := strconv.Atoi(message.Class)

		o.Room.Button = pos
		o.Room.Bet = 0
		o.Room.Cards = nil
		o.Room.Pot = make([]int, 1)
		o.Room.Each(0, func(o *poker.Occupant) bool {
			if o != nil {
				o.Bet = 0
				o.Action = ""
				o.Cards = nil
				o.Hand = 0
			}
			return true
		})

		dealer := o.Room.Occupants[pos-1]
		fmt.Printf("Button: %s(%s).\n", dealer.Id, dealer.Name)
	case poker.ActPreflop:
		fmt.Println("Preflop:", message.Class)
		cards := strings.Split(message.Class, ",")
		o.Cards = append(o.Cards, poker.ParseCard(cards[0]))
		o.Cards = append(o.Cards, poker.ParseCard(cards[1]))
	case poker.ActFlop:
		o.Room.Each(0, func(o *poker.Occupant) bool {
			if o != nil {
				o.Bet = 0
				o.Action = ""
			}
			return true
		})
		fmt.Println("Flop:", message.Class)
		cards := strings.Split(message.Class, ",")
		o.Room.Cards = append(o.Room.Cards, poker.ParseCard(cards[0]))
		o.Room.Cards = append(o.Room.Cards, poker.ParseCard(cards[1]))
		o.Room.Cards = append(o.Room.Cards, poker.ParseCard(cards[2]))
		o.Hand, _ = strconv.Atoi(cards[3])
	case poker.ActTurn:
		o.Room.Each(0, func(o *poker.Occupant) bool {
			if o != nil {
				o.Bet = 0
				o.Action = ""
			}
			return true
		})
		fmt.Println("Turn:", message.Class)
		cards := strings.Split(message.Class, ",")
		o.Room.Cards = append(o.Room.Cards, poker.ParseCard(cards[0]))
		o.Hand, _ = strconv.Atoi(cards[1])
	case poker.ActRiver:
		o.Room.Each(0, func(o *poker.Occupant) bool {
			if o != nil {
				o.Bet = 0
				o.Action = ""
			}
			return true
		})
		fmt.Println("River:", message.Class)
		cards := strings.Split(message.Class, ",")
		o.Room.Cards = append(o.Room.Cards, poker.ParseCard(cards[0]))
		o.Hand, _ = strconv.Atoi(cards[1])
	case poker.ActShowdown:
		fmt.Println("pot:", o.Room.Pot)
	case poker.ActAction:
		a := strings.Split(message.Class, ",")
		pos, _ := strconv.Atoi(a[0])
		o.Room.Bet, _ = strconv.Atoi(a[1])
		if o.Room.Occupants[pos-1].Id == o.Id {
			log.Printf("Your bet turn (%d/%d/%d):\n",
				o.Room.Occupants[pos-1].Bet, o.Room.Bet, o.Room.Occupants[pos-1].Chips)
		}
	case poker.ActPot:
		pots := strings.Split(message.Class, ",")
		o.Room.Pot = nil
		for i, _ := range pots {
			pot, _ := strconv.Atoi(pots[i])
			o.Room.Pot = append(o.Room.Pot, pot)
		}
	case poker.ActBet:
		occupant := o.Room.Occupant(message.From)
		occupant.Room = o.Room
		bets := strings.Split(message.Class, ",")
		occupant.Action = bets[0]
		occupant.Bet, _ = strconv.Atoi(bets[1])
		occupant.Chips, _ = strconv.Atoi(bets[2])

		if occupant.Id == o.Id {
			fmt.Printf("You %s: %d\n", occupant.Action, occupant.Bet)
		} else {
			fmt.Printf("%s(%s) %s: %d\n", occupant.Id, occupant.Name, occupant.Action, occupant.Bet)
		}
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
			if o.Room == nil {
				o.SendMessage(&poker.Message{
					Type:   poker.MsgPresence,
					Action: poker.ActJoin,
					To:     "0",
				})
			}

		case 'l':
			if o.Room != nil {
				o.SendMessage(&poker.Message{
					Type:   poker.MsgPresence,
					Action: poker.ActLeave,
					To:     "0",
				})
			}
			o.Pos = 0
			o.Room = nil
		case 'c':
			if o.Room != nil {
				cards := []poker.Card{}
				cards = append(cards, o.Cards...)

				cards = append(cards, o.Room.Cards...)
				fmt.Println(cards)
			}
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
