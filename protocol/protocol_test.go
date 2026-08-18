package protocol

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCreateGame(t *testing.T) {
	msg := CreateGame{Type: "CreateGame"}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	assertJSONKeys(t, data, map[string]interface{}{
		"type": "CreateGame",
	})

	var got CreateGame
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(msg, got) {
		t.Errorf("round-trip mismatch\n got: %+v\nwant: %+v", got, msg)
	}
}

func TestJoinGame(t *testing.T) {
	msg := JoinGame{Type: "JoinGame", GameID: "abc-123"}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	assertJSONKeys(t, data, map[string]interface{}{
		"type":   "JoinGame",
		"gameID": "abc-123",
	})

	var got JoinGame
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(msg, got) {
		t.Errorf("round-trip mismatch\n got: %+v\nwant: %+v", got, msg)
	}
}

func TestSubmitMove(t *testing.T) {
	msg := SubmitMove{
		Type:   "SubmitMove",
		GameID: "game-42",
		Target: Coord{Row: 3, Col: 7},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	assertJSONKeys(t, data, map[string]interface{}{
		"type":   "SubmitMove",
		"gameID": "game-42",
		"target": map[string]interface{}{"row": float64(3), "col": float64(7)},
	})

	var got SubmitMove
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(msg, got) {
		t.Errorf("round-trip mismatch\n got: %+v\nwant: %+v", got, msg)
	}
}

func TestStateUpdate(t *testing.T) {
	var board Board
	for c := 0; c < 5; c++ {
		board[0][c] = CellShip
	}

	msg := StateUpdate{
		Type:   "StateUpdate",
		GameID: "abc-123",
		State: GameState{
			MatchID: "game-abc",
			Status:  InProgress,
			YourView: PlayerView{
				Ships: []Ship{
					{
						Type:  Carrier,
						Cells: []Coord{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}},
						Sunk:  false,
					},
				},
				Board: board,
			},
			TheirView: PlayerView{
				Ships: []Ship{},
				Board: Board{},
			},
			YourTurn: true,
			Winner:   -1,
		},
		LastMove: LastMove{
			Coord:  Coord{Row: 3, Col: 7},
			Result: Hit,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}

	assertKey(t, raw, "type", "StateUpdate")
	assertKey(t, raw, "gameID", "abc-123")

	state, ok := raw["state"].(map[string]interface{})
	if !ok {
		t.Fatal("state is not a JSON object")
	}
	assertKey(t, state, "matchID", "game-abc")
	assertKey(t, state, "status", float64(InProgress))
	assertKey(t, state, "yourTurn", true)
	assertKey(t, state, "winner", float64(-1))

	yourView, ok := state["yourView"].(map[string]interface{})
	if !ok {
		t.Fatal("yourView is not a JSON object")
	}
	if _, ok := yourView["ships"]; !ok {
		t.Error("yourView.ships missing")
	}
	if _, ok := yourView["board"]; !ok {
		t.Error("yourView.board missing")
	}

	theirView, ok := state["theirView"].(map[string]interface{})
	if !ok {
		t.Fatal("theirView is not a JSON object")
	}
	if ships, ok := theirView["ships"].([]interface{}); !ok || len(ships) != 0 {
		t.Errorf("theirView.ships: got %v, want empty array", theirView["ships"])
	}

	lastMove, ok := raw["lastMove"].(map[string]interface{})
	if !ok {
		t.Fatal("lastMove is not a JSON object")
	}
	assertKey(t, lastMove, "result", float64(Hit))

	var got StateUpdate
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(msg, got) {
		t.Errorf("round-trip mismatch\n got: %+v\nwant: %+v", got, msg)
	}
}

func TestErrorResponse(t *testing.T) {
	msg := ErrorResponse{
		Type:    "ErrorResponse",
		Code:    400,
		Message: "not your turn",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	assertJSONKeys(t, data, map[string]interface{}{
		"type":    "ErrorResponse",
		"code":    float64(400),
		"message": "not your turn",
	})

	var got ErrorResponse
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(msg, got) {
		t.Errorf("round-trip mismatch\n got: %+v\nwant: %+v", got, msg)
	}
}

func TestTypeOf(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    string
		wantErr bool
	}{
		{"CreateGame", `{"type":"CreateGame"}`, "CreateGame", false},
		{"JoinGame", `{"type":"JoinGame","gameID":"x"}`, "JoinGame", false},
		{"SubmitMove", `{"type":"SubmitMove","gameID":"x","target":{"row":1,"col":2}}`, "SubmitMove", false},
		{"StateUpdate", `{"type":"StateUpdate","gameID":"x"}`, "StateUpdate", false},
		{"ErrorResponse", `{"type":"ErrorResponse","code":500,"message":"fail"}`, "ErrorResponse", false},
		{"invalid JSON", `{bad`, "", true},
		{"missing type", `{"gameID":"x"}`, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TypeOf([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Fatalf("err=%v, wantErr=%v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func assertJSONKeys(t *testing.T, data []byte, want map[string]interface{}) {
	t.Helper()
	var got map[string]interface{}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal for key check: %v", err)
	}
	for k, v := range want {
		assertKey(t, got, k, v)
	}
}

func assertKey(t *testing.T, m map[string]interface{}, key string, want interface{}) {
	t.Helper()
	got, ok := m[key]
	if !ok {
		t.Errorf("key %q missing from JSON", key)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("key %q: got %v (%T), want %v (%T)", key, got, got, want, want)
	}
}
