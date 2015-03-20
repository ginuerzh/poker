package main

import (
	"flag"
	"github.com/ginuerzh/poker/server"
	"log"
	"net/http"
)

var (
	WebDir string
	Addr   string
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	flag.StringVar(&Addr, "addr", ":8989", "server address ip:port")
	flag.StringVar(&WebDir, "web", "", "web directory rooted path")
	flag.Parse()
}

func main() {
	http.HandleFunc("/ws", poker.PokerHandler)
	http.Handle("/", http.FileServer(http.Dir(WebDir)))
	log.Fatal(http.ListenAndServe(Addr, nil))
}
