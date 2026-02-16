package Engine

import (
	"log"
	"net/url"
	"tworoomsapi/Logging"
)

// Given the path, gets any data that should be rendered with that requested template, if any. Only returns an error if one occurs (i.e. no data being found is not considered an error)
func GetApiData(path string, query url.Values) (any, error) {
	funcLogPrefix := "==GetApiData=="
	defer Logging.EnsureLogPrefixIsReset()
	Logging.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Getting Api Data for path {%s}", funcLogPrefix, path)

	// if strings.ToLower(path) == "/maps" {
	// 	mapIds, err := GetAllMaps()
	// 	return mapIds, err
	// } else if strings.ToLower(path) == "/recap" {
	// 	recap := GetRecap(query)
	// 	return recap, nil
	// }
	return nil, nil
}

// func GetRecap(query url.Values) Recap.Recap {
// 	funcLogPrefix := "==GetRecap=="
// 	roomCode := query.Get("roomCode")

// 	lobby, err := GetLobbyFromFs(roomCode)
// 	if err != nil {
// 		LogError(funcLogPrefix, err)
// 		return Recap.Recap{}
// 	}

// 	if lobby.Status != Models.LobbyStatus_Ended {
// 		return Recap.Recap{
// 			MapName: "Game has not ended yet",
// 		}
// 	}

// 	recap, err := Recap.GetRecapFromFs(lobby.GameStateId)
// 	if err != nil {
// 		LogError(funcLogPrefix, err)
// 		return Recap.Recap{}
// 	}
// 	return recap
// }
