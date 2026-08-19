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

var game_id string
var working bool

func MakeWork() {
	working = true
}
func IsWorking() bool {
	return working
}

func NewGame(c *websocket.Conn, err error, r *http.Request) {

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
