package riotinterface

import (
	"encoding/json"
	"les-randoms/utils"
)

type LeagueQueue string

const (
	SoloDuo LeagueQueue = "RANKED_SOLO_5x5"
	Flex    LeagueQueue = "RANKED_FLEX_SR"
)

type LeagueList struct {
	LeagueId string            `json:"leagueId"`
	Entries  []LeagueListEntry `json:"entries"`
	Tier     string            `json:"tier"`
	Name     string            `json:"name"`
	Queue    string            `json:"queue"`
}

type LeagueListEntry struct {
	FreshBlood   bool   `json:"freshBlood"`
	Wins         int    `json:"wins"`
	SummonerName string `json:"summonerName"`
	MiniSeries   []struct {
		Losses   int    `json:"losses"`
		Progress string `json:"progress"`
		Target   int    `json:"target"`
		Wins     int    `json:"wins"`
	} `json:"miniSeries"`
	Inactive     bool   `json:"inactive"`
	Veteran      bool   `json:"veteran"`
	HotStreak    bool   `json:"hotStreak"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Losses       int    `json:"losses"`
	SummonerId   string `json:"summonerId"`
}

func GetSoloDuoChallengerLeague() (*LeagueList, error) {
	body, err := getChallengerLeagueJSON(string(SoloDuo))
	if err != nil {
		return nil, err
	}
	league, err := parseLeagueListJSON(body)
	if err != nil {
		return nil, err
	}
	return league, nil
}

func getChallengerLeagueJSON(queue string) ([]byte, error) {
	return requestRIOTAPI("https://euw1.api.riotgames.com/lol/league/v4/challengerleagues/by-queue/" + queue)
}

func parseLeagueListJSON(body []byte) (*LeagueList, error) {
	data := &LeagueList{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	return data, nil
}
