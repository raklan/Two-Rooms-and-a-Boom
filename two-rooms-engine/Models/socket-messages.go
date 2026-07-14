package Models

// The different types of messages the server might send to a client connected via websocket.
const (
	WebsocketMessage_ClientStartGame        = "StartGame"
	WebsocketMessage_ClientStartRound       = "StartRound"
	WebsocketMessage_ClientNominateLeader   = "NominateLeader"
	WebsocketMessage_ClientAbdicate         = "Abdicate"
	WebsocketMessage_ClientAcceptAbdication = "AcceptAbdication"
	WebsocketMessage_ClientRejectAbdication = "RejectAbdication"
	WebsocketMessage_ClientUsurp            = "Usurp"
	WebsocketMessage_ClientUsurpVote        = "UsurpVote"
	WebsocketMessage_ClientHostageExchange  = "HostageExchange"
	WebsocketMessage_ClientCardShare        = "CardShare"
	WebsocketMessage_ClientAcceptCardShare  = "AcceptCardShare"

	WebsocketMessage_Close                   = "Close"
	WebsocketMessage_Error                   = "Error"
	WebsocketMessage_LobbyInfo               = "LobbyInfo"
	WebsocketMessage_GameState               = "GameState"
	WebsocketMessage_RoundStart              = "RoundStart"
	WebsocketMessage_RoundEnd                = "RoundEnd"
	WebsocketMessage_GameOver                = "GameOver"
	WebsocketMessage_NewLeader               = "NewLeader"
	WebsocketMessage_PendingAbdication       = "PendingAbdication"
	WebsocketMessage_AbdicationRejected      = "AbdicationRejected"
	WebsocketMessage_PendingUsurption        = "PendingUsurption"
	WebsocketMessage_UsurptionFailed         = "UsurptionFailed"
	WebsocketMessage_PendingCardShare        = "PendingCardShare"
	WebsocketMessage_CardShareRejected       = "CardShareRejected"
	WebsocketMessage_CardShare               = "CardShare"
	WebsocketMessage_HostageExchangeComplete = "HostageExchangeComplete"
)

type WebsocketMessageListItem struct {
	Message         WebsocketMessage
	ShouldBroadcast bool
}

// A message sent from the server to a client. The frontend can check [Type] to determine how to parse the object in [Data]
type WebsocketMessage struct {
	//One of the above constants. That constant will tell you which of the below structs is found in the [Data] field
	Type string `json:"type"`
	//One of the below structs, a Changelog, or a GameState. Its exact type is recorded in [Type]
	Data any `json:"data"`
}

// #region Client-sent Messages

type NominateLeader struct {
	//The ID of the player being nominated as the new leader
	NominatedPlayerId string `json:"nominatedPlayerId"`
}

type AcceptAbdication struct {
	Accept        bool   `json:"accept"`
	TakingOverFor string `json:"takingOverFor"`
}

type UsurpVote struct {
	//A bool indicating approval of the proposed new leader
	Vote bool `json:"vote"`
}

type HostageExchange struct {
	//An array of PlayerIDs indicating which players are to switch rooms
	Players []string `json:"players"`
}

type CardShareRequest struct {
	//Id of the player to share with
	ShareWith string `json:"shareWith"`
	//A boolean indicating whether both team and role should be shared (instead of just team)
	ShareFullCard bool `json:"shareFullCard"`
}

// #endregion

// #region Server-sent Messages

// A message containing a Player's assigned ID and the details of the lobby after they've joined it, whether by hosting it or joining a pre-existing lobby.
// The frontend should store this PlayerID.
type LobbyInfo struct {
	PlayerID  string `json:"playerID"`
	LobbyInfo Lobby  `json:"lobbyInfo"`
}

// A message to notify all players that a round has just started
type RoundStart struct {
	//The round number that is starting
	RoundNumber int `json:"roundNumber"`
	//The length of time, in seconds, for which the round will run
	RoundLength int `json:"roundLength"`
}

// A message to notify all players that a round has just ended
type RoundEnd struct {
}

type NewLeader struct {
	NewLeaderName string `json:"newLeaderName"`
}

type PendingAbdicationNotification struct {
	From string `json:"from"`
}

type PendingUsurptionNotification struct {
	UsurperName   string `json:"usurperName"`
	NewLeaderName string `json:"newLeaderName"`
}

type PendingCardShareNotification struct {
	FromName  string `json:"fromName"`
	FullShare bool   `json:"fullShare"`
}

type CardShare struct {
	FromPlayer string `json:"fromPlayer"`
	Team       string `json:"team"`
	Role       string `json:"role"`
}

type HostageExchangeResult struct {
	// map[room] -> list of players now in that room
	SentPlayers map[int][]string `json:"sentPlayers"`
}

// If some message from a client causes any error, one of these is sent back to the client
type SocketError struct {
	Message string `json:"message"`
}

// If a connection is about to be closed by the server, it will send a SocketClose, followed by immediately closing the connection
type SocketClose struct {
	Message string `json:"message"`
}

// A message informing a client the game has ended, containing a boolean describing whether they won
type GameOver struct {
	Victory bool `json:"victory"`
}

// #endregion
