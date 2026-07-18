package Models

import "slices"

type GameState struct {
	Id           string     `json:"id"`
	Players      []Player   `json:"players"`
	CurrentRound int        `json:"currentRound"`
	GameConfig   GameConfig `json:"gameConfig"`
}

func (g GameState) GetPlayersInRoom(room int) []Player {
	players := []Player{}

	for _, player := range g.Players {
		if player.Room == room {
			players = append(players, player)
		}
	}

	return players
}

func (g GameState) GetObscuredPlayersInRoom(room int) []PlayerObscured {
	players := []PlayerObscured{}

	for _, player := range g.Players {
		if player.Room == room {
			players = append(players, PlayerObscured{
				Id:           player.Id,
				Name:         player.Name,
				Room:         player.Room,
				IsRoomLeader: player.IsRoomLeader,
			})
		}
	}

	return players
}

func (g GameState) GetPlayerById(playerId string) (int, Player) {
	if playerIndex := slices.IndexFunc(g.Players, func(p Player) bool { return p.Id == playerId }); playerIndex == -1 {
		return -1, Player{}
	} else {
		return playerIndex, g.Players[playerIndex]
	}

}
