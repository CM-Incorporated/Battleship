package gamecode

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
  "cmincorporated.com/protocol"
)

var game_id string
var working bool

const maxAttemptsPerShip = 100

func PlaceShips(seed int64) (protocol.Board, []protocol.Ship) {
	rng := rand.New(rand.NewPCG(uint64(seed), uint64(seed)))

	for {
		var board protocol.Board
		ships := make([]protocol.Ship, 0, len(protocol.AllShipTypes))

		ok := true
		for _, t := range protocol.AllShipTypes {
			cells, placed := tryPlaceShip(&board, rng, t)
			if !placed {
				ok = false
				break
			}
			ships = append(ships, protocol.Ship{Type: t, Cells: cells})
		}
		if ok {
			return board, ships
		}
	}
}

func tryPlaceShip(board *protocol.Board, rng *rand.Rand, t protocol.ShipType) ([]protocol.Coord, bool) {
	length := t.Length()
	for range maxAttemptsPerShip {
		horizontal := rng.IntN(2) == 0
		var row, col int
		if horizontal {
			row = rng.IntN(protocol.GridSize)
			col = rng.IntN(protocol.GridSize - length + 1)
		} else {
			row = rng.IntN(protocol.GridSize - length + 1)
			col = rng.IntN(protocol.GridSize)
		}

		cells := shipCells(protocol.Coord{Row: row, Col: col}, length, horizontal)
		if !canPlace(board, cells) {
			continue
		}
		markShip(board, cells)
		return cells, true
	}
	return nil, false
}

func shipCells(start protocol.Coord, length int, horizontal bool) []protocol.Coord {
	cells := make([]protocol.Coord, length)
	for i := range length {
		if horizontal {
			cells[i] = protocol.Coord{Row: start.Row, Col: start.Col + i}
		} else {
			cells[i] = protocol.Coord{Row: start.Row + i, Col: start.Col}
		}
	}
	return cells
}

func canPlace(board *protocol.Board, cells []protocol.Coord) bool {
	for _, c := range cells {
		if !inBound(c) {
			return false
		}
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				n := protocol.Coord{Row: c.Row + dr, Col: c.Col + dc}
				if inBound(n) && board[n.Row][n.Col] != protocol.Water {
					return false
				}
			}
		}
	}
	return true
}

func markShip(board *protocol.Board, cells []protocol.Coord) {
	for _, c := range cells {
		board[c.Row][c.Col] = protocol.CellShip
	}
}

func inBound(c protocol.Coord) bool {
	return c.Row >= 0 && c.Row < protocol.GridSize && c.Col >= 0 && c.Col < protocol.GridSize
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
