package main

import (
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/ginuerzh/poker/server"
)

var (
	WebRoot   string
	Addr      string
	MongoAddr string
	RedisAddr string
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	flag.StringVar(&Addr, "addr", ":8989", "server address ip:port")
	flag.StringVar(&WebRoot, "web", "", "web directory rooted path")
	flag.Parse()
}

func main() {
	poker := &poker.Poker{
		Addr:    Addr,
		WebRoot: WebRoot,
		OnAuth:  onAuth,
	}
	log.Fatal(poker.ListenAndServe())
}

func onAuth(conn *poker.Conn, mechanism, text string) (*poker.Occupant, error) {
	id := strconv.FormatInt(time.Now().Unix(), 10)
	if text != "" {
		id = text
	}
	o := poker.NewOccupant(id, conn)
	o.Name = id

	return o, nil
}
