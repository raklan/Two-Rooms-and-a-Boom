package Models

import "sync"

type PendingAbdication struct {
	From string
	To   string
}

// map[roomCode][room#] -> playerId
var PendingAbdications = make(map[string]map[int]PendingAbdication)
var PendingAbdicationsMutex = sync.Mutex{}

type PendingUsurption struct {
	VotesYes           int
	VotesNo            int
	ProposedLeaderName string
}

// map[roomCode][room#] -> PendingUsurption
var PendingUsurptions = make(map[string]map[int]PendingUsurption)
var PendingUsurptionsMutex = sync.Mutex{}

type PendingCardShare struct {
	FromId        string
	ToId          string
	ShareFullCard bool
}

// map[roomCode][room#] -> []PendingCardShare
var PendingCardShares = make(map[string]map[int][]PendingCardShare)
var PendingCardSharesMutex = sync.Mutex{}
