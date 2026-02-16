package Engine

import (
	"fmt"
	"log"
	"math/rand"
	"tworoomsapi/Logging"
	"tworoomsengine/Models"
)

// Given an id to a Game defition, constructs and returns an initial GameState for it. This is essentially
// how to start the game
func GetInitialGameState(roomCode string, gameConfig Models.GameConfig) (Models.GameState, error) {
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

	assignTeams(&gameState)
	// if err := AssignRoles(&gameState, gameConfig.ActiveRoles, gameConfig.RequiredRoles); err != nil {
	// 	LogError(funcLogPrefix, err)
	// 	return gameState, err
	// }
	//TODO: Assign starting room

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

// func SubmitAction(gameId string, action Actions.SubmittedAction) ([]Models.WebsocketMessageListItem, error) {
// 	funcLogPrefix := "==SubmitAction=="
// 	defer Logging.EnsureLogPrefixIsReset()
// 	Logging.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

// 	gameState, err := GetGameStateFromFs(gameId)
// 	if err != nil {
// 		LogError(funcLogPrefix, err)
// 		return []Models.WebsocketMessageListItem{}, err
// 	}

// 	messages := []Models.WebsocketMessageListItem{}

// 	switch action.Type {
// 	case Actions.Action_Attack:
// 		turn := Actions.Attack{}
// 		err := json.Unmarshal(action.Turn, &turn)
// 		if err != nil {
// 			LogError(funcLogPrefix, err)
// 			return messages, err
// 		}

// 		var result any = nil
// 		if turn.IsAttacking() {
// 			result, err = turn.Execute(&gameState, action.PlayerId)
// 			messages = append(messages, Models.WebsocketMessageListItem{
// 				Message: Models.WebsocketMessage{
// 					Type: Models.WebsocketMessage_GameEvent,
// 					Data: result,
// 				},
// 				ShouldBroadcast: true,
// 			})
// 			messages = append(messages, Models.WebsocketMessageListItem{
// 				Message: Models.WebsocketMessage{
// 					Type: Models.WebsocketMessage_TurnEnd,
// 					Data: Models.TurnEnd{
// 						PlayerCurrentState: *gameState.GetCurrentPlayer(),
// 					},
// 				},
// 				ShouldBroadcast: false,
// 			})
// 		} else {
// 			// Cards are weird to deal with. You should only draw a card if you're in a dangerous sector,
// 			// so DrawCard will check where you are and set the type to Card_NoCard if so. In that case,
// 			// everyone should be told the player has moved into a safe sector, effectively skipping over
// 			// the Noise phase of their turn
// 			cardEvent, er := Actions.DrawCard(&gameState, action.PlayerId)
// 			err = er
// 			if cardEvent.Type == Models.Card_NoCard {

// 				actingPlayer := gameState.GetCurrentPlayer()

// 				if gameState.GameMap.Spaces[actingPlayer.GetSpaceMapKey()].Type != Models.Space_Pod {
// 					messages = append(messages, Models.WebsocketMessageListItem{
// 						Message: Models.WebsocketMessage{
// 							Type: Models.WebsocketMessage_GameEvent,
// 							Data: Models.GameEvent{
// 								Row:         -99,
// 								Col:         "!",
// 								Description: fmt.Sprintf("Player '%s' is in a safe sector", actingPlayer.Name),
// 							},
// 						},
// 						ShouldBroadcast: true,
// 					})
// 					go Recap.AddDataToRecap(gameId, action.PlayerId, gameState.Turn, "In a Safe Sector")
// 				}
// 				messages = append(messages, Models.WebsocketMessageListItem{
// 					Message: Models.WebsocketMessage{
// 						Type: Models.WebsocketMessage_TurnEnd,
// 						Data: Models.TurnEnd{
// 							PlayerCurrentState: *gameState.GetCurrentPlayer(),
// 						},
// 					},
// 					ShouldBroadcast: false,
// 				})
// 			} else {
// 				messages = append(messages, Models.WebsocketMessageListItem{
// 					Message: Models.WebsocketMessage{
// 						Type: Models.WebsocketMessage_Card,
// 						Data: cardEvent,
// 					},
// 					ShouldBroadcast: false,
// 				})
// 			}
// 		}

// 		if err != nil {
// 			LogError(funcLogPrefix, err)
// 			return messages, err
// 		}
// 	case Actions.Action_Movement:
// 		turn := Actions.Movement{}
// 		err := json.Unmarshal(action.Turn, &turn)
// 		if err != nil {
// 			LogError(funcLogPrefix, err)
// 			return messages, err
// 		}

// 		result, err := turn.Execute(&gameState, action.PlayerId)
// 		if err != nil {
// 			LogError(funcLogPrefix, err)
// 			return messages, err
// 		}

// 		messages = append(messages, Models.WebsocketMessageListItem{
// 			Message: Models.WebsocketMessage{
// 				Type: Models.WebsocketMessage_MovementResponse,
// 				Data: result,
// 			},
// 			ShouldBroadcast: false,
// 		})
// 	case Actions.Action_Noise:
// 		turn := Actions.Noise{}
// 		err := json.Unmarshal(action.Turn, &turn)
// 		if err != nil {
// 			LogError(funcLogPrefix, err)
// 			return messages, err
// 		}

// 		result, err := turn.Execute(&gameState, action.PlayerId)
// 		if err != nil {
// 			LogError(funcLogPrefix, err)
// 			return messages, err
// 		}

// 		messages = append(messages, Models.WebsocketMessageListItem{
// 			ShouldBroadcast: true,
// 			Message: Models.WebsocketMessage{
// 				Type: Models.WebsocketMessage_GameEvent,
// 				Data: result,
// 			},
// 		})

// 		messages = append(messages, Models.WebsocketMessageListItem{
// 			Message: Models.WebsocketMessage{
// 				Type: Models.WebsocketMessage_TurnEnd,
// 				Data: Models.TurnEnd{
// 					PlayerCurrentState: *gameState.GetCurrentPlayer(),
// 				},
// 			},
// 			ShouldBroadcast: false,
// 		})

// 	case Actions.Action_EndTurn:
// 		turn := Actions.EndTurn{}
// 		err := json.Unmarshal(action.Turn, &turn)
// 		if err != nil {
// 			LogError(funcLogPrefix, err)
// 			return messages, err
// 		}

// 		event, err := turn.Execute(&gameState, action.PlayerId)
// 		if err != nil {
// 			LogError(funcLogPrefix, err)
// 			return messages, err
// 		}

// 		if event != nil {
// 			messages = append(messages, Models.WebsocketMessageListItem{
// 				Message: Models.WebsocketMessage{
// 					Type: Models.WebsocketMessage_GameEvent,
// 					Data: event,
// 				},
// 				ShouldBroadcast: true,
// 			})
// 		}

// 		messages = append(messages, Models.WebsocketMessageListItem{
// 			Message: Models.WebsocketMessage{
// 				Type: Models.WebsocketMessage_GameState,
// 				Data: gameState,
// 			},
// 			ShouldBroadcast: true,
// 		})

// 	case Actions.Action_PlayCard:
// 		turn := Actions.PlayCard{}
// 		if err := json.Unmarshal(action.Turn, &turn); err != nil {
// 			LogError(funcLogPrefix, err)
// 			return messages, err
// 		}

// 		event, err := turn.Execute(&gameState, action.PlayerId)
// 		if err != nil {
// 			LogError(funcLogPrefix, err)
// 			return messages, err
// 		}

// 		messages = append(messages, Models.WebsocketMessageListItem{
// 			Message: Models.WebsocketMessage{
// 				Type: Models.WebsocketMessage_GameEvent,
// 				Data: event,
// 			},
// 			ShouldBroadcast: true,
// 		})

// 		messages = append(messages, Models.WebsocketMessageListItem{
// 			Message: Models.WebsocketMessage{
// 				Type: Models.WebsocketMessage_GameState,
// 				Data: gameState,
// 			},
// 			ShouldBroadcast: true,
// 		})
// 	}

// 	//Automatically end the game when there are no humans left
// 	numHumansLeft := 0
// 	for _, player := range gameState.Players {
// 		if player.Team == Models.PlayerTeam_Human {
// 			numHumansLeft++
// 		}
// 	}

// 	if numHumansLeft == 0 {
// 		messages = append(messages, Models.WebsocketMessageListItem{
// 			Message: Models.WebsocketMessage{
// 				Type: Models.WebsocketMessage_GameOver,
// 				Data: Models.GameOver{},
// 			},
// 			ShouldBroadcast: true,
// 		})
// 	}

// 	//Re-save gamestate
// 	_, err = SaveGameStateToFs(gameState)
// 	if err != nil {
// 		LogError(funcLogPrefix, err)
// 		return messages, err
// 	}

// 	return messages, nil
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

//#endregion
