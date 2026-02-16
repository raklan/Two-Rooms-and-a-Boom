package Models

const (
	LobbyStatus_AwaitingStart = "Awaiting Start"
	LobbyStatus_InProgress    = "In Progress"
	LobbyStatus_Ended         = "Game Ended"
)

type Lobby struct {
	RoomCode    string   `json:"roomCode"`
	Status      string   `json:"status"`
	Host        Player   `json:"host"`
	Players     []Player `json:"players"`
	MaxPlayers  int      `json:"maxPlayers"`
	GameStateId string   `json:"gameStateId"`
}
