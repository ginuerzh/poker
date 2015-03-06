package main

import (
	"github.com/ginuerzh/poker"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	http.HandleFunc("/", poker.PokerHandler)
	http.Handle("/poker", http.FileServer(http.Dir("/www/data/poker")))
	log.Fatal(http.ListenAndServe(":8989", nil))
}
