package gamecode

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func Connect(c *websocket.Conn, r *http.Request, channel chan string) {

	defer c.CloseNow()
	// Set the context as needed. Use of r.Context() is not recommended
	// to avoid surprising behavior (see http.Hijacker).
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var v any
	v = ""
	var err error

	log.Printf("New Player Connected")

	for v != "close" {
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			fmt.Println("Socket read error")
			return
		}
		log.Printf("received from %s: %v", r.RemoteAddr, v)

		// TODO: TALK TO GAME CODE FROM HERE
		channel <- v.(string)
	}

	c.Close(websocket.StatusNormalClosure, "")
}
