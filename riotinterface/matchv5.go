package riotinterface

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"les-randoms/utils"
	"net/http"
)

func GetMatchInfo(matchId string) ([]byte, error) {
	url := "https://europe.api.riotgames.com/lol/match/v5/matches/" + matchId
	utils.LogClassic("RIOT API REQUEST: " + url)

	request, error := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if error != nil {
		return nil, error
	}
	prepareRequestHeader(request)

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()

	if isResponseStatusOK(response.Status) {
		return nil, errors.New("RIOT API RESPONSE " + response.Status)
	}

	//fmt.Println("response Headers:", response.Header)

	return ioutil.ReadAll(response.Body)
}

func ParseMatchJSON(body []byte) (*Match, error) {
	data := &Match{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		utils.LogError(err.Error())
	}
	return data, nil
}
