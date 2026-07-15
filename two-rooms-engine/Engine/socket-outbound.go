package Engine

import (
	"fmt"
	"log"
	"math/rand"
	"slices"
	"time"
	"tworoomsapi/Logging"
	"tworoomsengine/Models"

	"github.com/gorilla/websocket"
)

// Given an id to a Game defition, constructs and returns an initial GameState for it. This is essentially
// how to start the game
func getInitialGameState(roomCode string, gameConfig Models.GameConfig) (Models.GameState, error) {
	funcLogPrefix := "==GetInitialGameState=="
	defer Logging.EnsureLogPrefixIsReset()
	Logging.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	gameState := Models.GameState{}

	lobby, err := GetLobbyFromFs(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	//Check if the lobby is already started.
	if lobby.Status == Models.LobbyStatus_InProgress {
		err := fmt.Errorf("tried to start game, but Lobby {%s} has been marked as In Progress and has a GameStateId == {%s}", lobby.RoomCode, lobby.GameStateId)
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	gameState.GameConfig = gameConfig

	gameState.Players = []Models.Player{}

	for _, element := range lobby.Players {
		gameState.Players = append(gameState.Players, Models.Player{
			Id:   element.Id,
			Name: element.Name,
			Team: element.Team,
			Role: "",
		})
	}

	gameState.CurrentRound = 1
	gameState.GameConfig.NumRounds = 3
	gameState.GameConfig.NumBlueTeam = 5 //TODO: don't hardcode
	gameState.GameConfig.NumRedTeam = 5

	assignTeams(&gameState)
	// if err := AssignRoles(&gameState, gameConfig.ActiveRoles, gameConfig.RequiredRoles); err != nil {
	// 	LogError(funcLogPrefix, err)
	// 	return gameState, err
	// }
	assignStartingRooms(&gameState)

	gameState, err = SaveGameStateToFs(gameState)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	//Mark the lobby as started and fill in GameStateId
	lobby.GameStateId = gameState.Id
	lobby.Status = Models.LobbyStatus_InProgress
	_, err = SaveLobbyToFs(lobby)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	//Create the recap object
	//go createInitialRecap(gameState)

	return gameState, nil
}

func EndGame(roomCode string, playerId string) error {
	funcLogPrefix := "==EndGame=="
	defer Logging.EnsureLogPrefixIsReset()
	Logging.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	lobby, err := GetLobbyFromFs(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	//Make sure that A) this player is the host and therefore allowed to end the game, and B) this game isn't already ended

	if lobby.Host.Id != playerId {
		return fmt.Errorf("player trying to end game is not host of lobby")
	}

	if lobby.Status == Models.LobbyStatus_Ended {
		return fmt.Errorf("game has already been marked as ended")
	}

	//Mark Game as ended and resave
	lobby.Status = Models.LobbyStatus_Ended

	_, err = SaveLobbyToFs(lobby)

	//Return any error that occurred during saving, if any
	return err
}

func MarkLobbyAsEnded(roomCode string) error {
	funcLogPrefix := "==MarkLobbyAsEnded=="
	defer Logging.EnsureLogPrefixIsReset()
	Logging.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	log.Printf("%s marking lobby as ended", funcLogPrefix)

	lobby, err := GetLobbyFromFs(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	if lobby.Status == Models.LobbyStatus_Ended {
		return fmt.Errorf("game has already been marked as ended")
	}

	//Mark Game as ended and resave
	lobby.Status = Models.LobbyStatus_Ended

	_, err = SaveLobbyToFs(lobby)

	//Return any error that occurred during saving, if any
	return err
}

// #region Round Starting/Ending
func startNextRound(roomCode string) error {
	funcLogPrefix := "==StartNextRound=="

	log.Printf("%s Starting next round for Room [%s]\n", funcLogPrefix, roomCode)

	lobby, err := GetLobbyFromFs(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	gameState, err := GetGameStateFromFs(lobby.GameStateId)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	log.Printf("%s determining round length for round %d\n", funcLogPrefix, gameState.CurrentRound)
	roundDurationSeconds := 0

	switch gameState.CurrentRound {
	case 1:
		roundDurationSeconds = 180
	case 2:
		roundDurationSeconds = 120
	case 3:
		roundDurationSeconds = 60
	}

	if roundDurationSeconds == 0 {
		return fmt.Errorf("could not determine round length")
	}
	log.Printf("%s round length will be %d seconds\n", funcLogPrefix, roundDurationSeconds)

	//TODO: DEBUG STATEMENT
	log.Printf("%s debug -- overriding round length to 10 seconds\n", funcLogPrefix)
	roundDurationSeconds = 10

	time.AfterFunc(time.Duration(roundDurationSeconds)*time.Second, func() {
		endCurrentRound(lobby.GameStateId, roomCode)
	})
	log.Printf("%s timer set to end round after %d seconds\n", funcLogPrefix, roundDurationSeconds)

	gamesClientsMutex.Lock()
	defer gamesClientsMutex.Unlock()

	for _, player := range gameState.Players {
		message := Models.WebsocketMessage{
			Type: Models.WebsocketMessage_RoundStart,
			Data: Models.RoundStart{
				RoundNumber: gameState.CurrentRound,
				RoundLength: roundDurationSeconds,
				Room:        player.Room,
			},
		}

		gamesClients[roomCode][player.Id].WriteJSON(message)
	}

	return nil
}

func endCurrentRound(gameStateId string, roomCode string) {
	funcLogPrefix := "==endCurrentRound=="
	log.Printf("%s ending Round for Room [%s]\n", funcLogPrefix, roomCode)

	gameState, err := GetGameStateFromFs(gameStateId)
	if err != nil {
		LogError(funcLogPrefix, err)
		return
	}

	gameState.CurrentRound++
	gamesClientsMutex.Lock()
	if gameState.CurrentRound <= gameState.GameConfig.NumRounds {
		log.Printf("%s dismissing leaders...", funcLogPrefix)
		for i := range gameState.Players {
			gameState.Players[i].IsRoomLeader = false
		}

		log.Printf("%s Current round is now %d/%d\n", funcLogPrefix, gameState.CurrentRound, gameState.GameConfig.NumRounds)
		sendMessageToAllPlayersInLobby(gamesClients[roomCode], Models.WebsocketMessage{
			Type: Models.WebsocketMessage_RoundEnd,
			Data: Models.RoundEnd{},
		})
	} else {
		log.Printf("%s round (%d) has surpassed max rounds (%d)\n", funcLogPrefix, gameState.CurrentRound, gameState.GameConfig.NumRounds)
		sendMessageToAllPlayersInLobby(gamesClients[roomCode], Models.WebsocketMessage{
			Type: Models.WebsocketMessage_GameOver,
			Data: Models.GameOver{},
		})
	}
	gamesClientsMutex.Unlock()

	_, err = SaveGameStateToFs(gameState)
	if err != nil {
		LogError(funcLogPrefix, err)
	}

}

// #region NominateLeader
func nominateLeader(roomCode string, nominatedPlayerId string) error {
	funcLogPrefix := "==nominateLeader=="

	log.Printf("%s Player Id [%s] has been nominated as a leader. Finding player in lobby...\n", funcLogPrefix, nominatedPlayerId)

	gameState, err := getGameStateFromRoomCode(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	nominatedPlayerIndex := slices.IndexFunc(gameState.Players, func(p Models.Player) bool { return p.Id == nominatedPlayerId })
	if nominatedPlayerIndex == -1 {
		LogError(funcLogPrefix, fmt.Errorf("could not find player with id %s to make room leader", nominatedPlayerId))
	}
	nominatedPlayer := gameState.Players[nominatedPlayerIndex]
	log.Printf("%s nominated player is Player [%s] in Room %d", funcLogPrefix, nominatedPlayer.Name, nominatedPlayer.Room)

	playersInRoom := gameState.GetPlayersInRoom(nominatedPlayer.Room)
	if i := slices.IndexFunc(playersInRoom, func(p Models.Player) bool { return p.IsRoomLeader }); i != -1 {
		err := fmt.Errorf("couldn't set player as leader - leader of room is already Player [%s]", playersInRoom[i].Name)
		LogError(funcLogPrefix, err)
		return err
	}

	gameState.Players[nominatedPlayerIndex].IsRoomLeader = true

	gamesClientsMutex.Lock()
	defer gamesClientsMutex.Unlock()

	sendMessageToAllPlayersInLobby(gamesClients[roomCode], Models.WebsocketMessage{
		Type: Models.WebsocketMessage_NewLeader,
		Data: Models.NewLeader{
			NewLeaderName: nominatedPlayer.Name,
		},
	})

	return nil
}

// #region Abdication
func abdicateLeadership(roomCode string, abdicatorId string, abdicateToId string) error {
	funcLogPrefix := "==abdicateLeadership=="

	gameState, err := getGameStateFromRoomCode(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	abdicatorIndex, abdicator := gameState.GetPlayerById(abdicateToId)
	if abdicatorIndex == -1 {
		err := fmt.Errorf("could not find player with Id == [%s]", abdicatorId)
		LogError(funcLogPrefix, err)
		return err
	}

	abdicateToIndex, abdicateTo := gameState.GetPlayerById(abdicateToId)
	if abdicateToIndex == -1 {
		err := fmt.Errorf("could not find player with Id == [%s]", abdicateToId)
		LogError(funcLogPrefix, err)
		return err
	}

	if !abdicator.IsRoomLeader {
		err := fmt.Errorf("player '%s' is not the room leader", abdicator.Name)
		LogError(funcLogPrefix, err)
		return err
	}

	if abdicator.Room != abdicateTo.Room {
		err := fmt.Errorf("cannot abdicate to Player '%s' -- Not in the same room", abdicateTo.Name)
		LogError(funcLogPrefix, err)
		return err
	}

	Models.PendingAbdicationsMutex.Lock()
	defer Models.PendingAbdicationsMutex.Unlock()

	if _, exists := Models.PendingAbdications[roomCode][abdicator.Room]; exists {
		err := fmt.Errorf("room already has a pending abdication from Player '%s' to Player '%s'", abdicator.Name, abdicateTo.Name)
		LogError(funcLogPrefix, err)
		return err
	}

	Models.PendingAbdications[roomCode][abdicator.Room] = Models.PendingAbdication{
		From: abdicator.Id,
		To:   abdicateTo.Id,
	}

	gamesClients[roomCode][abdicateToId].WriteJSON(Models.WebsocketMessage{
		Type: Models.WebsocketMessage_PendingAbdication,
		Data: Models.PendingAbdicationNotification{
			From: abdicator.Name,
		},
	})

	return nil
}

func respondAbdication(roomCode string, respondingPlayerId string, accept bool) error {
	funcLogPrefix := "==acceptAbdication=="

	log.Printf("%s accepting abdication for Player [%s]", funcLogPrefix, respondingPlayerId)

	gameState, err := getGameStateFromRoomCode(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	acceptorIndex, acceptor := gameState.GetPlayerById(respondingPlayerId)
	if acceptorIndex == -1 {
		err = fmt.Errorf("could not find player [%s]", respondingPlayerId)
		LogError(funcLogPrefix, err)
		return err
	}

	Models.PendingAbdicationsMutex.Lock()
	defer Models.PendingAbdicationsMutex.Unlock()

	if pendingAbdication, exists := Models.PendingAbdications[roomCode][acceptor.Room]; exists {
		if pendingAbdication.To != acceptor.Id {
			err = fmt.Errorf("player '%s' is not the target of the pending abdication in Room %d", acceptor.Name, acceptor.Room)
			LogError(funcLogPrefix, err)
			return err
		}

		abdicatorIndex, abdicator := gameState.GetPlayerById(pendingAbdication.From)
		if abdicatorIndex == -1 {
			err := fmt.Errorf("could not find abdicator with Id == [%s] in Pending Abdication to Player '%s'", pendingAbdication.From, acceptor.Name)
			LogError(funcLogPrefix, err)
			return err
		}

		gamesClientsMutex.Lock()
		defer gamesClientsMutex.Unlock()
		if accept {
			gameState.Players[abdicatorIndex].IsRoomLeader = false
			gameState.Players[acceptorIndex].IsRoomLeader = true

			sendMessageToAllPlayersInRoom(
				gamesClients[roomCode],
				gameState,
				acceptor.Room,
				Models.WebsocketMessage{
					Type: Models.WebsocketMessage_NewLeader,
					Data: Models.NewLeader{
						NewLeaderName: acceptor.Name,
					},
				})
		} else {
			gamesClients[roomCode][abdicator.Id].WriteJSON(Models.WebsocketMessage{
				Type: Models.WebsocketMessage_AbdicationRejected,
			})
		}

		delete(Models.PendingAbdications[roomCode], acceptor.Room)
	} else {
		err = fmt.Errorf("there is no pending abdication in Room %d", acceptor.Room)
		LogError(funcLogPrefix, err)
		return err
	}

	_, err = SaveGameStateToFs(gameState)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	return nil
}

// #region Usurption
func usurpLeader(roomCode string, usurperId string, proposedLeaderId string) error {
	funcLogPrefix := "==usurpLeader=="

	log.Printf("%s Player [%s] is trying to usurp room leader. Proposed new leader is Player [%s]", usurperId, proposedLeaderId)

	gameState, err := getGameStateFromRoomCode(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	usurperIndex, usurper := gameState.GetPlayerById(usurperId)
	if usurperIndex == -1 {
		err = fmt.Errorf("could not find usurper with Id == [%s]", usurperId)
		LogError(funcLogPrefix, err)
		return err
	}

	proposedLeaderIndex, proposedLeader := gameState.GetPlayerById(proposedLeaderId)
	if proposedLeaderIndex == -1 {
		err = fmt.Errorf("could not find proposed leader with Id == [%s]", proposedLeaderId)
		LogError(funcLogPrefix, err)
		return err
	}

	if usurper.Room != proposedLeader.Room {
		err = fmt.Errorf("new proposed leader is not in the same room as the usurper")
		LogError(funcLogPrefix, err)
		return err
	}

	Models.PendingUsurptionsMutex.Lock()
	defer Models.PendingAbdicationsMutex.Unlock()

	if _, exists := Models.PendingUsurptions[roomCode][usurper.Room]; exists {
		err = fmt.Errorf("there is already a pending usurption in Room %d", usurper.Room)
		LogError(funcLogPrefix, err)
		return err
	}

	Models.PendingUsurptions[roomCode][usurper.Room] = Models.PendingUsurption{
		VotesYes:           0,
		VotesNo:            0,
		ProposedLeaderName: proposedLeader.Name,
	}

	gamesClientsMutex.Lock()
	defer gamesClientsMutex.Unlock()

	sendMessageToAllPlayersInRoom(
		gamesClients[roomCode],
		gameState,
		usurper.Room,
		Models.WebsocketMessage{
			Type: Models.WebsocketMessage_PendingUsurption,
			Data: Models.PendingUsurptionNotification{
				UsurperName:   usurper.Name,
				NewLeaderName: proposedLeader.Name,
			},
		})

	return nil
}

func voteForUsurption(roomCode string, voterId string, vote bool) error {
	funcLogPrefix := "==voteForUsurption=="

	gameState, err := getGameStateFromRoomCode(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	voterIndex, voter := gameState.GetPlayerById(voterId)
	if voterIndex == -1 {
		err = fmt.Errorf("could not find player with Id == [%s]", voterId)
		LogError(funcLogPrefix, err)
		return err
	}

	Models.PendingUsurptionsMutex.Lock()
	defer Models.PendingAbdicationsMutex.Unlock()

	if currentVotes, exists := Models.PendingUsurptions[roomCode][voter.Room]; exists {
		if vote {
			Models.PendingUsurptions[roomCode][voter.Room] = Models.PendingUsurption{
				VotesYes:           currentVotes.VotesYes + 1,
				VotesNo:            currentVotes.VotesNo,
				ProposedLeaderName: currentVotes.ProposedLeaderName,
			}
		} else {
			Models.PendingUsurptions[roomCode][voter.Room] = Models.PendingUsurption{
				VotesYes:           currentVotes.VotesYes,
				VotesNo:            currentVotes.VotesNo + 1,
				ProposedLeaderName: currentVotes.ProposedLeaderName,
			}
		}

		newVotes := Models.PendingUsurptions[roomCode][voter.Room]
		numPlayersInRoom := len(gameState.GetPlayersInRoom(voter.Room))
		//If everyone has voted, end voting and notify of result
		if newVotes.VotesNo+newVotes.VotesYes == numPlayersInRoom {
			gamesClientsMutex.Lock()
			defer gamesClientsMutex.Unlock()
			var messageToSend = Models.WebsocketMessage{}
			if newVotes.VotesYes > numPlayersInRoom/2 {
				messageToSend = Models.WebsocketMessage{
					Type: Models.WebsocketMessage_NewLeader,
					Data: Models.NewLeader{
						NewLeaderName: newVotes.ProposedLeaderName,
					},
				}
			} else {
				messageToSend = Models.WebsocketMessage{
					Type: Models.WebsocketMessage_UsurptionFailed,
				}
			}

			sendMessageToAllPlayersInRoom(
				gamesClients[roomCode],
				gameState,
				voter.Room,
				messageToSend,
			)

			delete(Models.PendingUsurptions[roomCode], voter.Room)
		}
	} else {
		err = fmt.Errorf("there is no pending Usurption in Room %d", voter.Room)
		LogError(funcLogPrefix, err)
		return err
	}

	return nil
}

// #region CardShare
func requestCardShare(roomCode string, fromId string, toId string, fullShare bool) error {
	funcLogPrefix := "==requestCardShare=="

	log.Printf("%s received CardShare request from Player [%s] to Player [%s]\n", funcLogPrefix, fromId, toId)

	gameState, err := getGameStateFromRoomCode(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	fromPlayerIndex, fromPlayer := gameState.GetPlayerById(fromId)
	if fromPlayerIndex == -1 {
		err = fmt.Errorf("could not find Player [%s]", fromId)
		LogError(funcLogPrefix, err)
		return err
	}

	toPlayerIndex, toPlayer := gameState.GetPlayerById(toId)
	if toPlayerIndex == -1 {
		err = fmt.Errorf("could not find Player [%s]", toId)
		LogError(funcLogPrefix, err)
		return err
	}

	if fromPlayer.Room != toPlayer.Room {
		err = fmt.Errorf("can't cardshare with Player '%s' -- not in the same room!", toPlayer.Name)
		LogError(funcLogPrefix, err)
		return err
	}

	Models.PendingCardSharesMutex.Lock()
	defer Models.PendingCardSharesMutex.Unlock()

	roomCardShares := Models.PendingCardShares[roomCode][fromPlayer.Room]

	if slices.ContainsFunc(roomCardShares, func(el Models.PendingCardShare) bool { return el.FromId == fromPlayer.Id }) {
		err = fmt.Errorf("there is already an active Card Share request from Player '%s'", fromPlayer.Name)
		LogError(funcLogPrefix, err)
		return err
	}

	Models.PendingCardShares[roomCode][fromPlayer.Room] = append(roomCardShares, Models.PendingCardShare{
		FromId:        fromPlayer.Id,
		ToId:          toPlayer.Id,
		ShareFullCard: fullShare,
	})

	gamesClientsMutex.Lock()
	defer gamesClientsMutex.Unlock()

	gamesClients[roomCode][toPlayer.Id].WriteJSON(Models.WebsocketMessage{
		Type: Models.WebsocketMessage_PendingCardShare,
		Data: Models.PendingCardShareNotification{
			FromName:  fromPlayer.Name,
			FullShare: fullShare,
		},
	})

	return nil
}

func respondCardShare(roomCode string, respondingId string, accept bool) error {
	funcLogPrefix := "==responseCardShare=="

	log.Printf("%s responding to Card Share request to Player [%s]\n", funcLogPrefix, respondingId)

	gameState, err := getGameStateFromRoomCode(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	responderIndex, responder := gameState.GetPlayerById(respondingId)
	if responderIndex == -1 {
		err = fmt.Errorf("could not find Player with Id == [%s]", respondingId)
		LogError(funcLogPrefix, err)
		return err
	}

	Models.PendingCardSharesMutex.Lock()
	defer Models.PendingAbdicationsMutex.Unlock()

	pendingRoomShares := Models.PendingCardShares[roomCode][responder.Room]
	shareIndex := slices.IndexFunc(pendingRoomShares, func(el Models.PendingCardShare) bool { return el.ToId == responder.Id })
	if shareIndex == -1 {
		err = fmt.Errorf("no pending Card Share requests found to Player '%s'", responder.Name)
		LogError(funcLogPrefix, err)
		return err
	}

	shareRequest := pendingRoomShares[shareIndex]
	gamesClientsMutex.Lock()
	defer gamesClientsMutex.Unlock()
	if accept {
		requesterIndex, requester := gameState.GetPlayerById(shareRequest.FromId)
		if requesterIndex == -1 {
			err = fmt.Errorf("could not find Player with Id == [%s]", shareRequest.FromId)
			LogError(funcLogPrefix, err)
			return err
		}

		messageForRequester := Models.CardShare{
			FromPlayer: responder.Name,
			Team:       responder.Name,
		}
		messageForResponder := Models.CardShare{
			FromPlayer: requester.Name,
			Team:       requester.Name,
		}
		if shareRequest.ShareFullCard {
			messageForRequester.Role = responder.Role
			messageForResponder.Role = requester.Role
		}

		gamesClients[roomCode][requester.Id].WriteJSON(Models.WebsocketMessage{
			Type: Models.WebsocketMessage_CardShare,
			Data: messageForRequester,
		})
		gamesClients[roomCode][responder.Id].WriteJSON(Models.WebsocketMessage{
			Type: Models.WebsocketMessage_CardShare,
			Data: messageForResponder,
		})
	} else {
		gamesClients[roomCode][shareRequest.FromId].WriteJSON(Models.WebsocketMessage{
			Type: Models.WebsocketMessage_CardShareRejected,
		})
	}

	Models.PendingCardShares[roomCode][responder.Room] = slices.DeleteFunc(pendingRoomShares, func(el Models.PendingCardShare) bool { return el == shareRequest })

	return nil
}

// #region Hostage Exchange

func submitHostages(roomCode string, submittingPlayerId string, hostageIds []string) error {
	funcLogPrefix := "==submitHostages=="

	log.Printf("%s Submitting following hostages for Player [%s]: %s\n", funcLogPrefix, submittingPlayerId, hostageIds)

	gameState, err := getGameStateFromRoomCode(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	leaderIndex, leader := gameState.GetPlayerById(submittingPlayerId)
	if leaderIndex == -1 {
		err = fmt.Errorf("could not find player with ID == [%s]", submittingPlayerId)
		LogError(funcLogPrefix, err)
		return err
	}

	if !leader.IsRoomLeader {
		err = fmt.Errorf("Player '%s' is not the room leader and cannot submit hostages", leader.Name)
		LogError(funcLogPrefix, err)
		return err
	}

	finalHostageList := []string{}
	for _, id := range hostageIds {
		requestedHostageIndex, hostage := gameState.GetPlayerById(id)
		if requestedHostageIndex == -1 {
			log.Printf("%s could not find Hostage with Id == [%s]. Skipping...\n", funcLogPrefix, id)
			continue
		}

		if hostage.Room != leader.Room {
			log.Printf("%s Hostage is not in the same room as the Leader. Skipping...\n", funcLogPrefix)
			continue
		}

		finalHostageList = append(finalHostageList, id)
	}

	Models.PendingHostagesMutex.Lock()
	defer Models.PendingHostagesMutex.Unlock()

	if pendingExchange, exists := Models.PendingHostages[roomCode][leader.Room]; exists {
		hostagesToSendToThisRoom := pendingExchange.HostageIds
		hostagesToSendToOtherRoom := finalHostageList

		sentToThisRoom := []string{}
		sentToOtherRoom := []string{}
		for _, id := range hostagesToSendToThisRoom {
			hostageIndex, _ := gameState.GetPlayerById(id)
			if hostageIndex == -1 {
				log.Printf("%s could not find Hostage to send to this room with Id == [%s]. Skipping...\n", funcLogPrefix, id)
				continue
			}

			gameState.Players[hostageIndex].Room = leader.Room
			sentToThisRoom = append(sentToThisRoom, gameState.Players[hostageIndex].Name)
		}

		for _, id := range hostagesToSendToOtherRoom {
			hostageIndex, _ := gameState.GetPlayerById(id)
			if hostageIndex == -1 {
				log.Printf("%s could not find Hostage to send to other room with Id == [%s]. Skipping...\n", funcLogPrefix, id)
				continue
			}

			gameState.Players[hostageIndex].Room = pendingExchange.Room
			sentToOtherRoom = append(sentToOtherRoom, gameState.Players[hostageIndex].Name)
		}

		gamesClientsMutex.Lock()
		defer gamesClientsMutex.Unlock()
		sendMessageToAllPlayersInLobby(gamesClients[roomCode], Models.WebsocketMessage{
			Type: Models.WebsocketMessage_HostageExchangeComplete,
			Data: Models.HostageExchangeResult{
				SentPlayers: map[int][]string{
					pendingExchange.Room: sentToOtherRoom,
					leader.Room:          sentToThisRoom,
				},
			},
		})

		_, err = SaveGameStateToFs(gameState)
		if err != nil {
			LogError(funcLogPrefix, err)
			return err
		}
	} else {
		Models.PendingHostages[roomCode][leader.Room] = Models.PendingHostageExchange{
			Room:       leader.Room,
			HostageIds: finalHostageList,
		}
	}

	return nil
}

// #region Spectators

//Come back to this if I ever want spectators
// func SwitchPlayerSpectating(roomCode string, playerId string, isSpectating bool) (Models.Lobby, error) {
// 	funcLogPrefix := "==SwitchPlayerSpectating=="
// 	defer Logging.EnsureLogPrefixIsReset()
// 	Logging.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

// 	lobby, err := GetLobbyFromFs(roomCode)
// 	if err != nil {
// 		LogError(funcLogPrefix, err)
// 		return Models.Lobby{}, err
// 	}

// 	playerIndexToSwitch := slices.IndexFunc(lobby.Players, func(p Models.Player) bool { return p.Id == playerId })
// 	if playerIndexToSwitch == -1 {
// 		err := fmt.Errorf("could not find requested player with ID == %s", playerId)
// 		LogError(funcLogPrefix, err)
// 		return Models.Lobby{}, err
// 	}

// 	if isSpectating {
// 		lobby.Players[playerIndexToSwitch].Team = Models.PlayerTeam_Spectator
// 	} else {
// 		lobby.Players[playerIndexToSwitch].Team = ""
// 	}

// 	_, err = SaveLobbyToFs(lobby)

// 	return lobby, err
// }

// #region Helper Functions

// Assigns teams randomly to all players in the GameState. If a player cannot be assigned for any reason, they are assigned as a spectator
func assignTeams(gameState *Models.GameState) {
	log.Println("Assigning teams")
	blueToAssign, redToAssign := gameState.GameConfig.NumBlueTeam, gameState.GameConfig.NumRedTeam
	for index := range gameState.Players {
		// if gameState.Players[index].Team == Models.PlayerTeam_Spectator {
		// 	continue
		// }
		if blueToAssign == 0 && redToAssign != 0 { //No blue slots left, must be red
			gameState.Players[index].Team = Models.PlayerTeam_Red
		} else if redToAssign == 0 && blueToAssign != 0 { //No red slots left, must be blue
			gameState.Players[index].Team = Models.PlayerTeam_Blue
		} else {
			if blueToAssign == 0 && redToAssign == 0 {
				gameState.Players[index].Team = Models.PlayerTeam_Neutral
			}
			if rand.Intn(2) == 0 {
				gameState.Players[index].Team = Models.PlayerTeam_Blue
				blueToAssign = blueToAssign - 1
			} else {
				gameState.Players[index].Team = Models.PlayerTeam_Red
				redToAssign = redToAssign - 1
			}
		}
	}
}

func assignStartingRooms(gameState *Models.GameState) {
	log.Println("Assigning starting rooms")

	slices.SortFunc(gameState.Players, func(p1 Models.Player, p2 Models.Player) int { return rand.Intn(100) - rand.Intn(100) })

	for i := range gameState.Players {
		if i <= len(gameState.Players)/2 {
			gameState.Players[i].Room = 1
		} else {
			gameState.Players[i].Room = 2
		}
	}

	log.Println("Starting rooms assigned. Note - player list ordering has been changed.")
}

// func AssignRoles(gameState *Models.GameState, activeRoles map[string]int, requiredRoles map[string]int) error {
// 	log.Println("Assigning roles")
// 	humanPlayers := gameState.GetHumanPlayers()
// 	alienPlayers := gameState.GetAlienPlayers()

// 	for (len(humanPlayers) > 0 || len(alienPlayers) > 0) && len(requiredRoles) > 0 {
// 		maps.DeleteFunc(requiredRoles, func(name string, num int) bool { return num == 0 })
// 		for roleName := range requiredRoles {
// 			var playerListToAssignFrom []Models.Player

// 			if Models.RoleTeams[roleName] == Models.PlayerTeam_Human {
// 				playerListToAssignFrom = humanPlayers
// 			} else if Models.RoleTeams[roleName] == Models.PlayerTeam_Alien {
// 				playerListToAssignFrom = alienPlayers
// 			}

// 			if len(playerListToAssignFrom) == 0 {
// 				return fmt.Errorf("too many required roles for Team %s", Models.RoleTeams[roleName])
// 			}

// 			playerToAssign_Copy := playerListToAssignFrom[rand.Intn(len(playerListToAssignFrom))]

// 			if Models.RoleAssigners[roleName] == nil {
// 				LogError("AssignRoles", fmt.Errorf("ryan forgot to add the assigner for %s to the RoleAssigners", roleName))
// 			}

// 			Models.RoleAssigners[roleName](&gameState.Players[slices.IndexFunc(gameState.Players, func(p Models.Player) bool { return p.Id == playerToAssign_Copy.Id })])

// 			if Models.RoleTeams[roleName] == Models.PlayerTeam_Human {
// 				humanPlayers = slices.DeleteFunc(playerListToAssignFrom, func(p Models.Player) bool { return p.Id == playerToAssign_Copy.Id })
// 			} else if Models.RoleTeams[roleName] == Models.PlayerTeam_Alien {
// 				alienPlayers = slices.DeleteFunc(playerListToAssignFrom, func(p Models.Player) bool { return p.Id == playerToAssign_Copy.Id })
// 			}

// 			requiredRoles[roleName]--
// 			if requiredRoles[roleName] <= 0 {
// 				delete(requiredRoles, roleName)
// 			}
// 			if _, exists := activeRoles[roleName]; exists {
// 				activeRoles[roleName]--
// 			}
// 		}
// 	}

// 	for (len(humanPlayers) > 0 || len(alienPlayers) > 0) && len(activeRoles) > 0 {
// 		maps.DeleteFunc(activeRoles, func(name string, num int) bool { return num == 0 })
// 		if len(activeRoles) == 0 {
// 			break
// 		}
// 		roleName, _ := Models.GetRandomMapPair(activeRoles)

// 		var playerListToAssignFrom []Models.Player

// 		if Models.RoleTeams[roleName] == Models.PlayerTeam_Human {
// 			playerListToAssignFrom = humanPlayers
// 		} else if Models.RoleTeams[roleName] == Models.PlayerTeam_Alien {
// 			playerListToAssignFrom = alienPlayers
// 		}

// 		if len(playerListToAssignFrom) == 0 { //If we can't assign this role, delete it so we don't keep checking it
// 			delete(activeRoles, roleName)
// 			continue
// 		}

// 		playerToAssign_Copy := playerListToAssignFrom[rand.Intn(len(playerListToAssignFrom))]

// 		Models.RoleAssigners[roleName](&gameState.Players[slices.IndexFunc(gameState.Players, func(p Models.Player) bool { return p.Id == playerToAssign_Copy.Id })])

// 		if Models.RoleTeams[roleName] == Models.PlayerTeam_Human {
// 			humanPlayers = slices.DeleteFunc(playerListToAssignFrom, func(p Models.Player) bool { return p.Id == playerToAssign_Copy.Id })
// 		} else if Models.RoleTeams[roleName] == Models.PlayerTeam_Alien {
// 			alienPlayers = slices.DeleteFunc(playerListToAssignFrom, func(p Models.Player) bool { return p.Id == playerToAssign_Copy.Id })
// 		}

// 		activeRoles[roleName]--
// 		if activeRoles[roleName] <= 0 {
// 			delete(activeRoles, roleName)
// 		}
// 	}

// 	return nil
// }

func sendMessageToAllPlayersInLobby(lobby map[string]*websocket.Conn, message Models.WebsocketMessage) {
	funcLogPrefix := "==sendMessageToAllPlayersInLobby=="

	if message.Type == "" {
		log.Printf("%s WARNING: Websocket message being sent has no Type set! Frontend will likely not know how to handle the message!", funcLogPrefix)
	}

	for playerId, conn := range lobby {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Printf("%s Error sending message, skipping meesage to PlayerId [%s]", funcLogPrefix, playerId)
			continue
		}
	}
}

func sendMessageToAllPlayersInRoom(lobby map[string]*websocket.Conn, gameState Models.GameState, roomNumber int, message Models.WebsocketMessage) {
	funcLogPrefix := "==sendMessageToAllPlayersInRoom=="

	if message.Type == "" {
		log.Printf("%s WARNING: Websocket message being sent has no Type set! Frontend will likely not know how to handle the message!", funcLogPrefix)
	}

	playersInRoom := gameState.GetPlayersInRoom(roomNumber)

	for _, player := range playersInRoom {
		err := lobby[player.Id].WriteJSON(message)
		if err != nil {
			log.Printf("%s error sending message, skipping message to Player [%s]", funcLogPrefix, player.Id)
			continue
		}
	}
}

func cleanUpRoom(room map[string]*websocket.Conn, roomCode string) {
	funcLogPrefix := "==cleanUpRoom=="
	log.Printf("%s cleaning up Room {%s}", funcLogPrefix, roomCode)
	closeMessage := Models.WebsocketMessage{
		Type: Models.WebsocketMessage_Close,
		Data: Models.SocketClose{
			Message: "Game has ended. Closing connection",
		},
	}

	gamesClientsMutex.Lock()
	//Send the messages to every player and stop tracking their connection
	for playerId, conn := range room {
		log.Printf("%s stopping tracking and closing connection for PlayerId %s", funcLogPrefix, playerId)
		err := conn.WriteJSON(closeMessage)
		if err != nil {
			log.Printf("Error sending Close Message to %s. Aborting message, but closing connection anyways", playerId)
		}
		conn.Close()
		delete(room, playerId)
	}

	//Stop tracking the room
	log.Printf("%s stopping tracking of Room {%s}", funcLogPrefix, roomCode)
	delete(gamesClients, roomCode)
	gamesClientsMutex.Unlock()
	log.Printf("%s room successfully cleaned up", funcLogPrefix)
}

func endPlayerConnection(roomCode string, playerId string, room map[string]*websocket.Conn) (Models.Lobby, error) {
	funcLogPrefix := "==endPlayerConnection=="
	//Tell the engine to remove the player from the DB copy of the lobby
	updatedLobby, err := LeaveRoom(roomCode, playerId)
	if err != nil {
		LogError(funcLogPrefix, err)
		return Models.Lobby{}, err
	}

	//If the removed client has a currently open connection, tell that the client that the connection is closing, then close connection
	if conn, exists := room[playerId]; exists {
		msg := Models.WebsocketMessage{
			Type: Models.WebsocketMessage_Close,
			Data: Models.SocketClose{
				Message: "Player has been removed from Lobby. Closing connection",
			},
		}
		conn.WriteJSON(msg)
		conn.Close()
		//Remove connection from lobby map so we don't try to send them any more messages
		gamesClientsMutex.Lock()
		delete(room, playerId)
		gamesClientsMutex.Unlock()
	}

	return updatedLobby, nil
}

func sendError(conn *websocket.Conn, err error) {
	conn.WriteJSON(Models.WebsocketMessage{
		Type: Models.WebsocketMessage_Error,
		Data: Models.SocketError{
			Message: err.Error(),
		},
	})
}

func getGameStateFromRoomCode(roomCode string) (Models.GameState, error) {
	funcLogPrefix := "==getGameStateFromRoomCode=="

	lobby, err := GetLobbyFromFs(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return Models.GameState{}, err
	}

	gameState, err := GetGameStateFromFs(lobby.GameStateId)
	if err != nil {
		LogError(funcLogPrefix, err)
		return Models.GameState{}, err
	}

	return gameState, nil
}

//#endregion
