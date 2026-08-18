package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "cmincorporated.com/protocol"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println("Handler error")
	}
	defer c.CloseNow()

	// Set the context as needed. Use of r.Context() is not recommended
	// to avoid surprising behavior (see http.Hijacker).
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var v any
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		fmt.Println("First read error")
		return
	}
	log.Printf("received: %v", v)

	for v != "close" {
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			fmt.Println("Main loop error")
			return
		}
		log.Printf("received from %s: %v", r.RemoteAddr, v)
	}

	c.Close(websocket.StatusNormalClosure, "")
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
