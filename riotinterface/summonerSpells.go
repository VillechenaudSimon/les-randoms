package riotinterface

import (
	"encoding/json"
	"les-randoms/utils"
	"strconv"
)

var summonerSpellsMap map[string]SummonerSpell

type SummonerSpellsInfo struct { // Only 'useful' informations are parsed from JSON
	Type    string                   `json:"type"`
	Version string                   `json:"version"`
	Data    map[string]SummonerSpell `json:"data"`
}

type SummonerSpell struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	Image       riotImage `json:"image"`
}

func GetSummonerSpellsMap() map[string]SummonerSpell {
	updateServerInfoIfNecessary()
	return summonerSpellsMap
}

func GetSummonerSpellImageNameByKey(key int) string {
	for _, summonerSpell := range GetSummonerSpellsMap() {
		if summonerSpell.Key == strconv.Itoa(key) {
			return summonerSpell.Image.Full
		}
	}
	utils.LogError("Summoner Spell Key Not Found : " + strconv.Itoa(key))
	return ""
}

func getSummonerSpellsInfo() (*SummonerSpellsInfo, error) {
	body, err := getAllSummonerSpellsJSON()
	if err != nil {
		return nil, err
	}
	summonerSpellsInfo, err := parseSummonerSpellsJSON(body)
	if err != nil {
		return nil, err
	}
	return summonerSpellsInfo, nil
}

func getAllSummonerSpellsJSON() ([]byte, error) {
	return requestRIOTAPI("https://ddragon.leagueoflegends.com/cdn/11.23.1/data/en_US/summoner.json")
}

func parseSummonerSpellsJSON(body []byte) (*SummonerSpellsInfo, error) {
	data := &SummonerSpellsInfo{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	return data, nil
}

func updateServerSummonerSpellsInfo() error {
	summonerSpellsInfo, err := getSummonerSpellsInfo()
	if err != nil {
		utils.LogError("Error while updating server summoner spells info :\n" + err.Error())
		return err
	}
	summonerSpellsMap = summonerSpellsInfo.Data
	return nil
}
