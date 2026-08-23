package gamecode

import (
	"fmt"

	"cmincorporated.com/protocol"
)

func PlaceShips(seed int64) (protocol.Board, []protocol.Ship) {
	// Create board, place all 5 ships at random
	var board protocol.Board = xxx
}

func canPlace(protocol.Board, protocol.Coord) bool {
	// check if a ship can be placed at a specific location
	fmt.Printf("Yo")
	return true
}

func checkShipPlacementValid(protocol.Coord) bool {
	return true
}
