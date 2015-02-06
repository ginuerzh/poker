package main

import (
	"github.com/ginuerzh/poker"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", poker.PokerHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
