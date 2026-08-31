package main

import (
	"context"
	"log"
	"time"

	_ "cmincorporated.com/protocol"

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

	for i := 0; i < 10; i++ {

		err = wsjson.Write(ctx, c, "hello")
		if err != nil {
			log.Fatal("Oh god")
		}

		time.Sleep(time.Second)

	}

	err = wsjson.Write(ctx, c, "close")
	if err != nil {
		log.Fatal("Oh god")
	}

	c.Close(websocket.StatusNormalClosure, "")
}
