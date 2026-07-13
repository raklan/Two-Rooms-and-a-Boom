package Models

type GameConfig struct {
	NumBlueTeam    int            `json:"numBlueTeam"`
	NumRedTeam     int            `json:"numRedTeam"`
	NumNeutralTeam int            `json:"numNeutralTeam"`
	NumRounds      int            `json:"numRounds"`
	ActiveRoles    map[string]int `json:"activeRoles"`
	RequiredRoles  map[string]int `json:"requiredRoles"`
}
