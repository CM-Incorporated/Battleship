package protocol

import "encoding/json"

// Grid size constant.
const GridSize = 10

// CellState represents the state of a single cell on the board.
type CellState int

const (
	Water CellState = iota
	CellShip
	Hit
	Miss
	Sunk
)

// ShipType identifies a ship in the fleet.
type ShipType int

const (
	Carrier ShipType = iota
	Battleship
	Destroyer
	Cruiser
	Submarine
)

// GameStatus tracks the lifecycle of a match.
type GameStatus int

const (
	Waiting GameStatus = iota
	InProgress
	Finished
)

// Coord is a row/column position on the 10x10 grid.
type Coord struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// Board is a 10x10 grid of cell states.
type Board [GridSize][GridSize]CellState

// Ship records the type, position, and sink status of a single ship.
type Ship struct {
	Type  ShipType `json:"type"`
	Cells []Coord  `json:"cells"`
	Sunk  bool     `json:"sunk"`
}

// PlayerView is one player's perspective: their ships and their board.
type PlayerView struct {
	Ships []Ship `json:"ships"`
	Board Board  `json:"board"`
}

// GameState is the full state sent to a client in a StateUpdate.
type GameState struct {
	MatchID   string     `json:"matchID"`
	Status    GameStatus `json:"status"`
	YourView  PlayerView `json:"yourView"`
	TheirView PlayerView `json:"theirView"`
	YourTurn  bool       `json:"yourTurn"`
	Winner    int        `json:"winner"`
}

// LastMove describes the result of the most recent attack.
type LastMove struct {
	Coord  Coord     `json:"coord"`
	Result CellState `json:"result"`
}

// CreateGame asks the server to start a new game.
type CreateGame struct {
	Type string `json:"type"`
}

// JoinGame pairs a second player into an existing game.
type JoinGame struct {
	Type   string `json:"type"`
	GameID string `json:"gameID"`
}

// SubmitMove fires at a coordinate on the opponent's board.
type SubmitMove struct {
	Type   string `json:"type"`
	GameID string `json:"gameID"`
	Target Coord  `json:"target"`
}

// StateUpdate pushes the current game state to both clients.
type StateUpdate struct {
	Type     string    `json:"type"`
	GameID   string    `json:"gameID"`
	State    GameState `json:"state"`
	LastMove LastMove  `json:"lastMove"`
}

// ErrorResponse is sent when the server rejects a message.
type ErrorResponse struct {
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// TypeOf extracts the "type" discriminator from a raw JSON message
// without decoding the full payload. Returns the type string.
func TypeOf(data []byte) (string, error) {
	var wrapper struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return "", err
	}
	return wrapper.Type, nil
}
