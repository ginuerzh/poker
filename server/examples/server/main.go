package main

import (
	"flag"
	"github.com/ginuerzh/poker/server"
	"log"
	"strings"
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
	flag.StringVar(&MongoAddr, "mongo", "localhost:27017", "mongodb addr")
	flag.StringVar(&RedisAddr, "redis", "localhost:6379", "redis addr")
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

func onAuth(conn *poker.Conn, mechanism, token string) (*poker.Occupant, error) {
	id := onlineUser(token)
	if len(id) == 0 {
		return nil, poker.NewError(2, "user not found")
	}
	user := &User{}
	if err := user.FindById(id); err != nil {
		return nil, poker.NewError(3, "db error")
	}

	o := poker.NewOccupant(id, conn)

	o.Name = user.Nickname
	o.Profile = strings.Replace(user.Profile, "172.24.222.54", "172.24.222.42", -1)
	o.Chips = int(user.Chips)
	o.Level = int(user.Level())

	return o, nil
}
