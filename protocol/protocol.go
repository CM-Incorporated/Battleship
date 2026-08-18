package protocol

import "encoding/json"

const GridSize = 10

type CellState int

const (
	Water CellState = iota
	CellShip
	Hit
	Miss
	Sunk
)

type ShipType int

const (
	Carrier ShipType = iota
	Battleship
	Destroyer
	Cruiser
	Submarine
)

type GameStatus int

const (
	Waiting GameStatus = iota
	InProgress
	Finished
)

type Coord struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type Board [GridSize][GridSize]CellState

type Ship struct {
	Type  ShipType `json:"type"`
	Cells []Coord  `json:"cells"`
	Sunk  bool     `json:"sunk"`
}

type PlayerView struct {
	Ships []Ship `json:"ships"`
	Board Board  `json:"board"`
}

type GameState struct {
	MatchID   string     `json:"matchID"`
	Status    GameStatus `json:"status"`
	YourView  PlayerView `json:"yourView"`
	TheirView PlayerView `json:"theirView"`
	YourTurn  bool       `json:"yourTurn"`
	Winner    int        `json:"winner"`
}

type LastMove struct {
	Coord  Coord     `json:"coord"`
	Result CellState `json:"result"`
}

type CreateGame struct {
	Type string `json:"type"`
}

type JoinGame struct {
	Type   string `json:"type"`
	GameID string `json:"gameID"`
}

type SubmitMove struct {
	Type   string `json:"type"`
	GameID string `json:"gameID"`
	Target Coord  `json:"target"`
}

type StateUpdate struct {
	Type     string    `json:"type"`
	GameID   string    `json:"gameID"`
	State    GameState `json:"state"`
	LastMove LastMove  `json:"lastMove"`
}

type ErrorResponse struct {
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func TypeOf(data []byte) (string, error) {
	var wrapper struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return "", err
	}
	return wrapper.Type, nil
}
