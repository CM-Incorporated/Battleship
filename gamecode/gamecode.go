package gamecode

var game_id string
var working bool

func MakeWork() {
	working = true
}
func IsWorking() bool {
	return working
}
