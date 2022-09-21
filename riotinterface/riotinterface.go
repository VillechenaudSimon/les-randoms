package riotinterface

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"les-randoms/utils"
	"net/http"
	"os"
	"strings"
	"time"
)

var LastServerUpdatePatch string      // example : 12.17.1
var LastServerUpdatePatchShort string // example : 12.17
var VersionsList []string

var LastServerUpdateTime time.Time

type RiotApiError int
type GameModeId int

const (
	RiotApiErrorUnknown         RiotApiError = 0
	RiotApiErrorForbidden       RiotApiError = 403
	RiotApiErrorTooManyRequests RiotApiError = 429

	TeamBlueId int = 100
	TeamRedId  int = 200

	GameModeIdNormal GameModeId = 400
	GameModeIdSoloQ  GameModeId = 420
	GameModeIdFlexQ  GameModeId = 440
	GameModeIdARAM   GameModeId = 450
	GameModeIdClash  GameModeId = 700
)

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
	tmp := strings.Split(LastServerUpdatePatch, ".")
	LastServerUpdatePatchShort = tmp[0] + "." + tmp[1]
	utils.LogNotNilError(updateServerSummonerSpellsInfo())
	utils.LogNotNilError(updateServerItemsInfo())
	utils.LogNotNilError(updateServerRunesInfo())

	utils.LogInfo("Server informations updated")
}

func updateServerInfoIfNecessary() error {
	if LastServerUpdateTime.Add(time.Hour * 24).Before(time.Now()) {
		LastServerUpdateTime = time.Now()
		var err error
		VersionsList, err = GetVersionsArray()
		lastRiotVersion := VersionsList[0] //GetLastVersion()
		if err == nil {
			if LastServerUpdatePatch == "NOTUPDATEDYET" || LastServerUpdatePatch != lastRiotVersion {
				updateServerInfo(lastRiotVersion)
			}
		} else {
			return err
		}
	}
	return nil
}

func GetPatch() string {
	updateServerInfoIfNecessary()
	return LastServerUpdatePatch
}

func GetPatchShort() string {
	updateServerInfoIfNecessary()
	return LastServerUpdatePatchShort
}
