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
	log.Fatal(http.ListenAndServe(":8000", nil))
}
