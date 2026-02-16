package Engine

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"tworoomsapi/Logging"
	"tworoomsengine/Models"
)

func PrepareFilesystem() {
	os.Mkdir("./roles", 0666)
	os.Mkdir("./lobbies", 0666)
	os.Mkdir("./gameStates", 0666)
	os.Mkdir("./recaps", 0666)
}

// Yes, I know. I just REALLY didn't want to bring in an entire database JUST for this and Redis shouldn't be used for it
func SaveRoleToDB(m Models.Role) (Models.Role, error) {
	funcLogPrefix := "==SaveMapToDB=="

	asJson, err := json.Marshal(m)
	if err != nil {
		LogError(funcLogPrefix, err)
		return m, err
	}

	filename := "role_" + m.Name + ".json"
	f, err := os.Create(fmt.Sprintf("./roles/%s", filename))
	if err != nil {
		LogError(funcLogPrefix, err)
		f.Close()
		return m, err
	}
	_, err = f.Write(asJson)
	f.Close()
	if err != nil {
		LogError(funcLogPrefix, err)
		return m, err
	}

	return m, nil
}

func GetRoleFromDB(roleName string) (Models.Role, error) {
	funcLogPrefix := "==GetMapFromDB=="
	defer Logging.EnsureLogPrefixIsReset()
	Logging.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Getting role from DB with name == {%s}", funcLogPrefix, roleName)
	data, err := os.ReadFile(fmt.Sprintf("./roles/role_%s.json", roleName))
	if err != nil {
		LogError(funcLogPrefix, err)
		return Models.Role{}, err
	}

	parsed := Models.Role{}

	err = json.Unmarshal(data, &parsed)
	if err != nil {
		LogError(funcLogPrefix, err)
		return parsed, err
	}

	return parsed, nil
}

func SaveLobbyToFs(lobby Models.Lobby) (Models.Lobby, error) {
	funcLogPrefix := "==SaveLobbyToFs=="
	defer Logging.EnsureLogPrefixIsReset()
	Logging.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	asJson, err := json.Marshal(lobby)
	if err != nil {
		LogError(funcLogPrefix, err)
		return lobby, err
	}

	filename := "lobby_" + lobby.RoomCode + ".json"

	f, err := os.Create(fmt.Sprintf("./lobbies/%s", filename))
	if err != nil {
		LogError(funcLogPrefix, err)
		f.Close()
		return lobby, err
	}
	_, err = f.Write(asJson)
	f.Close()
	if err != nil {
		LogError(funcLogPrefix, err)
		return lobby, err
	}

	//Kick off goroutine clearing out unused lobbies
	go clearOutOldFiles("./lobbies/")

	return lobby, nil
}

func GetLobbyFromFs(roomCode string) (Models.Lobby, error) {
	funcLogPrefix := "==GetLobbyFromFs=="
	defer Logging.EnsureLogPrefixIsReset()
	Logging.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Getting lobby from FS with RoomCode == {%s}", funcLogPrefix, roomCode)
	data, err := os.ReadFile(fmt.Sprintf("./lobbies/lobby_%s.json", roomCode))
	if err != nil {
		LogError(funcLogPrefix, err)
		return Models.Lobby{}, err
	}

	parsed := Models.Lobby{}

	err = json.Unmarshal(data, &parsed)
	if err != nil {
		LogError(funcLogPrefix, err)
		return parsed, err
	}

	return parsed, nil
}

func SaveGameStateToFs(gameState Models.GameState) (Models.GameState, error) {
	funcLogPrefix := "==SaveGameStateToFs==:"
	defer Logging.EnsureLogPrefixIsReset()
	Logging.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Received GameState to save", funcLogPrefix)

	//If the gameState doesn't have an ID yet,
	//Generate one for it by simply using the Current UNIX time in milliseconds
	id := gameState.Id
	if id == "" {
		log.Printf("%s GameState does not yet have an ID. Generating new one.", funcLogPrefix)
		id = GenerateId()
		log.Printf("%s ID successfully generated. Assigning ID {%s} to GameState", funcLogPrefix, id)
		gameState.Id = id
	}

	asJson, err := json.Marshal(gameState)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	filename := "gameState_" + gameState.Id + ".json"
	f, err := os.Create(fmt.Sprintf("./gameStates/%s", filename))
	if err != nil {
		LogError(funcLogPrefix, err)
		f.Close()
		return gameState, err
	}
	_, err = f.Write(asJson)
	f.Close()
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	//Kick off goroutine clearing out unused lobbies
	go clearOutOldFiles("./gameStates/")
	return gameState, nil
}

func GetGameStateFromFs(id string) (Models.GameState, error) {
	funcLogPrefix := "==GetGameStateFromFs=="
	defer Logging.EnsureLogPrefixIsReset()
	Logging.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Getting GameState from FS with ID == {%s}", funcLogPrefix, id)
	data, err := os.ReadFile(fmt.Sprintf("./gameStates/gameState_%s.json", id))
	if err != nil {
		LogError(funcLogPrefix, err)
		return Models.GameState{}, err
	}

	parsed := Models.GameState{}

	err = json.Unmarshal(data, &parsed)
	if err != nil {
		LogError(funcLogPrefix, err)
		return parsed, err
	}

	return parsed, nil
}

func clearOutOldFiles(directory string) {
	files, _ := os.ReadDir(directory)
	for _, file := range files {
		fullFileName := directory + file.Name()
		stats, _ := os.Stat(fullFileName)
		expirationTime := stats.ModTime().AddDate(0, 0, 7)
		if time.Now().After(expirationTime) {
			os.Remove(fullFileName)
		}
	}
}
