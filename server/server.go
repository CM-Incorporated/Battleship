package main

import (
	"fmt"
	"log"
	"net/http"

	"cmincorporated.com/gamecode"
	_ "cmincorporated.com/protocol"
	"github.com/coder/websocket"
)

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println("Handler error")
	}
	go gamecode.NewGame(c, err, r)

}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
