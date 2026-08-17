package main

import (
	"context"
	"log"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	c, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.CloseNow()

	for i := 0; i < 100; i++ {
		err = wsjson.Write(ctx, c, i)
		if err != nil {
			log.Fatal("Oh god")
		}

	}

	err = wsjson.Write(ctx, c, "close")
	if err != nil {
		log.Fatal("Oh god")
	}

	c.Close(websocket.StatusNormalClosure, "")
}
