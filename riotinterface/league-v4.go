package riotinterface

import (
	"encoding/json"
	"errors"
	"les-randoms/utils"
	"strings"
)

type LeagueQueue string

const (
	SoloDuo LeagueQueue = "RANKED_SOLO_5x5"
	Flex    LeagueQueue = "RANKED_FLEX_SR"
)

var LeagueV4Errors leagueV4Errors

func init() {
	LeagueV4Errors = leagueV4Errors{
		MissingSummonerSoloDuoEntry: "No solo duo league entry found for this summoner",
	}
}

type leagueV4Errors struct {
	MissingSummonerSoloDuoEntry string
}

type LeagueListDTO struct {
	LeagueId string          `json:"leagueId"`
	Entries  []LeagueItemDTO `json:"entries"`
	Tier     string          `json:"tier"`
	Name     string          `json:"name"`
	Queue    string          `json:"queue"`
}

type LeagueItemDTO struct {
	FreshBlood   bool          `json:"freshBlood"`
	Wins         int           `json:"wins"`
	SummonerName string        `json:"summonerName"`
	MiniSeries   MiniSeriesDTO `json:"miniSeries"`
	Inactive     bool          `json:"inactive"`
	Veteran      bool          `json:"veteran"`
	HotStreak    bool          `json:"hotStreak"`
	Rank         string        `json:"rank"`
	LeaguePoints int           `json:"leaguePoints"`
	Losses       int           `json:"losses"`
	SummonerId   string        `json:"summonerId"`
}

type LeagueEntryDTO struct {
	LeagueId     string        `json:"leagueId"`
	SummonerId   string        `json:"summonerId"`
	SummonerName string        `json:"summonerName"`
	QueueType    string        `json:"queueType"`
	Tier         string        `json:"tier"`
	Rank         string        `json:"rank"`
	LeaguePoints int           `json:"leaguePoints"`
	Wins         int           `json:"wins"`
	Losses       int           `json:"losses"`
	HotStreak    bool          `json:"hotStreak"`
	Veteran      bool          `json:"veteran"`
	FreshBlood   bool          `json:"freshBlood"`
	Inactive     bool          `json:"inactive"`
	MiniSeries   MiniSeriesDTO `json:"miniSeries"`
}

type MiniSeriesDTO struct {
	Losses   int    `json:"losses"`
	Progress string `json:"progress"`
	Target   int    `json:"target"`
	Wins     int    `json:"wins"`
}

func ParseTierRank(tier string, rank string) string {
	switch tier {
	case "MASTER":
		return "Master"
	case "GRANDMASTER":
		return "GrandMaster"
	default:
		return strings.Title(strings.ToLower(tier)) + " " + rank
	}
}

func GetSoloDuoChallengerLeague() (*LeagueListDTO, error) {
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

func GetEntriesBySummonerId(id string) ([]*LeagueEntryDTO, error) {
	body, err := getEntriesBySummonerIdJSON(id)
	if err != nil {
		return nil, err
	}
	entries, err := parseEntriesJSON(body)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func GetSoloDuoEntryBySummonerId(id string) (*LeagueEntryDTO, error) {
	entries, err := GetEntriesBySummonerId(id)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.QueueType == string(SoloDuo) {
			return entry, nil
		}
	}
	return nil, errors.New(LeagueV4Errors.MissingSummonerSoloDuoEntry)
}

func getChallengerLeagueJSON(queue string) ([]byte, error) {
	return requestRIOTAPI("https://euw1.api.riotgames.com/lol/league/v4/challengerleagues/by-queue/" + queue)
}

func getEntriesBySummonerIdJSON(id string) ([]byte, error) {
	return requestRIOTAPI("https://euw1.api.riotgames.com/lol/league/v4/entries/by-summoner/" + id)
}

func parseLeagueListJSON(body []byte) (*LeagueListDTO, error) {
	data := &LeagueListDTO{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	return data, nil
}

func parseEntriesJSON(body []byte) ([]*LeagueEntryDTO, error) {
	data := make([]*LeagueEntryDTO, 0)
	err := json.Unmarshal(body, &data)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	return data, nil
}
