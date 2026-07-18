package Models

const (
	PlayerTeam_Blue    = "Blue"
	PlayerTeam_Red     = "Red"
	PlayerTeam_Neutral = "Neutral"
)

type Player struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Team         string `json:"team"`
	Role         string `json:"role"`
	Room         int    `json:"room"`
	IsRoomLeader bool   `json:"isRoomLeader"`
}

// Version of the Player struct without Team/Role info. Safe to send to any client
type PlayerObscured struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Room         int    `json:"room"`
	IsRoomLeader bool   `json:"isRoomLeader"`
}
