package main

import (
	"fmt"
	"log"
	"net/http"

	"cmincorporated.com/gamecode"
	_ "cmincorporated.com/protocol"
	"github.com/coder/websocket"
)

func lobby(queue chan connection) {
	var firstPlayer *connection

	for incomingPlayer := range queue {
		if firstPlayer == nil {
			p := incomingPlayer
			firstPlayer = &p
			fmt.Println("First player waiting...")
		} else {
			p1 := *firstPlayer
			p2 := incomingPlayer

			firstPlayer = nil

			fmt.Println("Match found! Launching game.")
			go gamecode.Connect(p1.c, p1.r, p1.channel)
			go gamecode.Connect(p2.c, p2.r, p2.channel)
			go gamecode.NewGame(p1.channel, p2.channel)
		}
	}
}

type connection struct {
	c       *websocket.Conn
	r       *http.Request
	channel chan string
}

func makeHandler(queue chan connection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			fmt.Println("Handler error")
			return
		}

		p := connection{c, r, make(chan string)}
		queue <- p
	}
}

func main() {
	lobbyQueue := make(chan connection)
	go lobby(lobbyQueue)

	http.HandleFunc("/", makeHandler(lobbyQueue))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
