package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"cmincorporated.com/gamecode"
	"cmincorporated.com/protocol"
)

func main() {
	seed := parseSeed()
	board, ships := gamecode.PlaceShips(seed)

	fmt.Printf("Seed: %d\n", seed)
	printBoard(board, ships)

	fmt.Println("\nShips:")
	for _, s := range ships {
		fmt.Printf("  %-11s len=%d sunk=%v cells=%v\n", shipName(s.Type), s.Type.Length(), s.Sunk, s.Cells)
	}
}

func parseSeed() int64 {
	if len(os.Args) > 1 {
		v, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err == nil {
			return v
		}
		fmt.Printf("invalid seed %q, using random\n", os.Args[1])
	}
	return time.Now().UnixNano()
}

func shipName(t protocol.ShipType) string {
	switch t {
	case protocol.Carrier:
		return "Carrier"
	case protocol.Battleship:
		return "Battleship"
	case protocol.Destroyer:
		return "Destroyer"
	case protocol.Cruiser:
		return "Cruiser"
	case protocol.Submarine:
		return "Submarine"
	}
	return "Unknown"
}

func shipLetter(t protocol.ShipType) byte {
	switch t {
	case protocol.Carrier:
		return 'A'
	case protocol.Battleship:
		return 'B'
	case protocol.Destroyer:
		return 'D'
	case protocol.Cruiser:
		return 'C'
	case protocol.Submarine:
		return 'S'
	}
	return '?'
}

func printBoard(board protocol.Board, ships []protocol.Ship) {
	mark := map[protocol.Coord]byte{}
	for _, s := range ships {
		for _, c := range s.Cells {
			mark[c] = shipLetter(s.Type)
		}
	}

	fmt.Print("\n     ")
	for c := range protocol.GridSize {
		fmt.Printf("%d ", c)
	}
	fmt.Println()

	for r := range protocol.GridSize {
		fmt.Printf("%2d | ", r)
		for c := range protocol.GridSize {
			coord := protocol.Coord{Row: r, Col: c}
			if ch, ok := mark[coord]; ok {
				fmt.Printf("%c ", ch)
				continue
			}
			switch board[r][c] {
			case protocol.Hit:
				fmt.Print("X ")
			case protocol.Miss:
				fmt.Print("o ")
			case protocol.Sunk:
				fmt.Print("@ ")
			case protocol.CellShip:
				fmt.Print("# ")
			default:
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}

	fmt.Println("\nLegend: A=Carrier B=Battleship D=Destroyer C=Cruiser S=Submarine")
}
