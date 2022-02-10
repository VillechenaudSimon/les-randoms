package riotinterface

import (
	"bytes"
	"errors"
	"fmt"
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
	return status[0:3] != fmt.Sprint(http.StatusOK) // 200 in theory
}

func requestRIOTAPI(url string) ([]byte, error) {
	utils.LogClassic("REQUEST RIOT API : " + url)

	request, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if err != nil {
		//utils.LogError("Error while creating request for riotAPI : " + err.Error())
		return nil, err
	}
	request.Header.Set("Accept-Charset", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("X-Riot-Token", os.Getenv("X_RIOT_TOKEN"))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		//utils.LogError("Error while requesting riotAPI : " + err.Error())
		return nil, err
	}
	defer response.Body.Close()

	if isResponseStatusNotOK(response.Status) {
		err = errors.New("RESPONSE RIOT API : " + response.Status)
		//utils.LogError("Error with riotapi response : " + err.Error())
		return nil, err
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

func GetPatch() string {
	updateServerInfoIfNecessary()
	return LastServerUpdatePatch
}

func ParseGameModeFromQueueId(id int) string {
	switch id {
	case -1: //Unknown code for ARAM
		return "ARAM"
	case 400:
		return "Normal Game"
	case 420:
		return "Ranked Solo/Duo"
	default:
		return "Unknown Game Mode (queueId : " + fmt.Sprint(id) + ")"
	}
}
