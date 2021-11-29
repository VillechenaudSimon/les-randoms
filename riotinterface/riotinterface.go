package riotinterface

import (
	"bytes"
	"errors"
	"io/ioutil"
	"les-randoms/utils"
	"net/http"
	"os"
	"time"
)

var LastServerUpdatePatch string
var LastServerUpdateTime time.Time

func init() {
	LastServerUpdatePatch = "NOTUPDATEDYET"
	updateServerInfoIfNecessary()
}

func isResponseStatusNotOK(status string) bool {
	return status[0:3] != "200"
}

func requestRIOTAPI(url string) ([]byte, error) {
	utils.LogClassic("REQUEST RIOT API : " + url)

	request, error := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if error != nil {
		return nil, error
	}
	request.Header.Set("Accept-Charset", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("X-Riot-Token", os.Getenv("X_RIOT_TOKEN"))

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()

	if isResponseStatusNotOK(response.Status) {
		return nil, errors.New("RESPONSE RIOT API : " + response.Status)
	}

	//fmt.Println("response Headers:", response.Header)

	return ioutil.ReadAll(response.Body)
}

// Call this function with newPatch="" if you want to request RIOT API for the last patch name automatically
func updateServerInfo(newPatch string) {
	if newPatch == "" {
		lastRiotVersion, err := GetLastVersion()
		if err != nil {
			return
		}
		LastServerUpdatePatch = lastRiotVersion
	} else {
		LastServerUpdatePatch = newPatch
	}
	utils.LogNotNilError(updateServerSummonerSpellsInfo())
	utils.LogNotNilError(updateServerItemsInfo())
}

func updateServerInfoIfNecessary() {
	if LastServerUpdateTime.Add(time.Hour * 24).Before(time.Now()) {
		LastServerUpdateTime = time.Now()
		lastRiotVersion, err := GetLastVersion()
		if err == nil {
			if LastServerUpdatePatch == "NOTUPDATEDYET" || LastServerUpdatePatch != lastRiotVersion {
				updateServerInfo(lastRiotVersion)
			}
		}
	}
}

func ParseGameMode(input string) string {
	switch input {
	case "ARAM":
		return "ARAM"
	case "CLASSIC":
		return "Ranked Solo/Duo"
	default:
		return "Unknown Game Mode"
	}
}

type riotImage struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
}
