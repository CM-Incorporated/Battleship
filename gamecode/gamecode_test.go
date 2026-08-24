package gamecode

import (
	"reflect"
	"testing"

	"cmincorporated.com/protocol"
)

func TestPlaceShipsDeterministic(t *testing.T) {
	b1, s1 := PlaceShips(42)
	b2, s2 := PlaceShips(42)
	if !reflect.DeepEqual(b1, b2) || !reflect.DeepEqual(s1, s2) {
		t.Fatal("same seed produced different placements")
	}
}

func TestPlaceShipsFleetIntegrity(t *testing.T) {
	totalLen := 0
	for _, st := range protocol.AllShipTypes {
		totalLen += st.Length()
	}

	for seed := int64(0); seed < 500; seed++ {
		board, ships := PlaceShips(seed)

		if len(ships) != len(protocol.AllShipTypes) {
			t.Fatalf("seed %d: got %d ships, want %d", seed, len(ships), len(protocol.AllShipTypes))
		}

		seen := make(map[protocol.Coord]bool)
		shipCellCount := 0
		for i, s := range ships {
			wantType := protocol.AllShipTypes[i]
			if s.Type != wantType {
				t.Fatalf("seed %d ship %d: type %v, want %v", seed, i, s.Type, wantType)
			}
			if len(s.Cells) != s.Type.Length() {
				t.Fatalf("seed %d ship %v: %d cells, want %d", seed, s.Type, len(s.Cells), s.Type.Length())
			}
			if s.Sunk {
				t.Fatalf("seed %d ship %v: fresh ship must not be sunk", seed, s.Type)
			}
			for j, c := range s.Cells {
				if !inBound(c) {
					t.Fatalf("seed %d ship %v: cell %v out of bounds", seed, s.Type, c)
				}
				if seen[c] {
					t.Fatalf("seed %d: duplicate cell %v across ships", seed, c)
				}
				seen[c] = true
				shipCellCount++
				if j > 0 {
					prev := s.Cells[j-1]
					straight := (prev.Row == c.Row && abs(prev.Col-c.Col) == 1) ||
						(prev.Col == c.Col && abs(prev.Row-c.Row) == 1)
					if !straight {
						t.Fatalf("seed %d ship %v: cells %v->%v not contiguous/straight", seed, s.Type, prev, c)
					}
				}
				if board[c.Row][c.Col] != protocol.CellShip {
					t.Fatalf("seed %d: board[%d][%d] not CellShip", seed, c.Row, c.Col)
				}
			}
		}
		if shipCellCount != totalLen {
			t.Fatalf("seed %d: %d ship cells, want %d", seed, shipCellCount, totalLen)
		}
		waterCount := 0
		for r := range protocol.GridSize {
			for c := range protocol.GridSize {
				if board[r][c] == protocol.Water {
					waterCount++
				}
			}
		}
		if waterCount != protocol.GridSize*protocol.GridSize-totalLen {
			t.Fatalf("seed %d: wrong water count %d", seed, waterCount)
		}
	}
}

func TestPlaceShipsNoTouching(t *testing.T) {
	for seed := int64(0); seed < 500; seed++ {
		_, ships := PlaceShips(seed)
		for i := range ships {
			for j := i + 1; j < len(ships); j++ {
				for _, a := range ships[i].Cells {
					for _, b := range ships[j].Cells {
						dr := a.Row - b.Row
						dc := a.Col - b.Col
						if dr >= -1 && dr <= 1 && dc >= -1 && dc <= 1 {
							t.Fatalf("seed %d: ships %v and %v touch at %v/%v", seed, ships[i].Type, ships[j].Type, a, b)
						}
					}
				}
			}
		}
	}
}

func TestDifferentSeedsVary(t *testing.T) {
	first, _ := PlaceShips(1)
	differs := false
	for seed := int64(2); seed < 50 && !differs; seed++ {
		b, _ := PlaceShips(seed)
		if !reflect.DeepEqual(first, b) {
			differs = true
		}
	}
	if !differs {
		t.Fatal("all seeds produced identical boards")
	}
}

func TestCanPlace(t *testing.T) {
	var board protocol.Board
	markShip(&board, []protocol.Coord{{Row: 5, Col: 5}, {Row: 5, Col: 6}})

	tests := []struct {
		name  string
		cells []protocol.Coord
		want  bool
	}{
		{"empty area", []protocol.Coord{{Row: 0, Col: 0}}, true},
		{"overlap", []protocol.Coord{{Row: 5, Col: 6}}, false},
		{"orthogonal touch", []protocol.Coord{{Row: 5, Col: 7}}, false},
		{"vertical touch", []protocol.Coord{{Row: 6, Col: 5}}, false},
		{"diagonal touch", []protocol.Coord{{Row: 6, Col: 7}}, false},
		{"gap ok", []protocol.Coord{{Row: 5, Col: 8}}, true},
		{"out of bounds", []protocol.Coord{{Row: -1, Col: 0}}, false},
		{"out of bounds high", []protocol.Coord{{Row: protocol.GridSize, Col: 0}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canPlace(&board, tt.cells); got != tt.want {
				t.Errorf("canPlace(%v) = %v, want %v", tt.cells, got, tt.want)
			}
		})
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
