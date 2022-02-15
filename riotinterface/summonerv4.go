package riotinterface

import (
	"encoding/json"
	"les-randoms/utils"
)

type Summoner struct {
	Id            string `json:"id"`
	AccountId     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"` // According to Riot Documentation this can be a long (int64?) but since a player needs 135 years of game time while leveluping once per hour to reach the maximum level an int can store, i think it's safe
}

func GetSummonerFromName(name string) (*Summoner, error) {
	body, err := getSummonerFromNameJSON(name)
	if err != nil {
		return nil, err
	}
	summoner, err := parseSummonerJSON(body)
	if err != nil {
		return nil, err
	}
	return summoner, nil
}

func GetPuuidFromSummonerName(name string) (string, error) {
	summoner, err := GetSummonerFromName(name)
	if err != nil {
		return "", err
	}
	return summoner.Puuid, nil
}

func getSummonerFromNameJSON(name string) ([]byte, error) {
	return requestRIOTAPI("https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + name)
}

func GetSummonerFromId(id string) (*Summoner, error) {
	body, err := getSummonerFromIdJSON(id)
	if err != nil {
		return nil, err
	}
	summoner, err := parseSummonerJSON(body)
	if err != nil {
		return nil, err
	}
	return summoner, nil
}

func getSummonerFromIdJSON(id string) ([]byte, error) {
	return requestRIOTAPI("https://euw1.api.riotgames.com/lol/summoner/v4/summoners/" + id)
}

func parseSummonerJSON(body []byte) (*Summoner, error) {
	data := &Summoner{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	return data, nil
}
