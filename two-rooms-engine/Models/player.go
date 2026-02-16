package Models

const (
	PlayerTeam_Blue    = "Blue"
	PlayerTeam_Red     = "Red"
	PlayerTeam_Neutral = "Neutral"
)

type Player struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Team string `json:"team"`
	Role string `json:"role"`
}
