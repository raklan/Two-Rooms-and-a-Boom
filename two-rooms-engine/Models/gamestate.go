package Models

type GameState struct {
	Id           string     `json:"id"`
	Players      []Player   `json:"players"`
	CurrentRound int        `json:"currentRound"`
	GameConfig   GameConfig `json:"gameConfig"`
}
