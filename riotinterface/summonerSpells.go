package riotinterface

import (
	"encoding/json"
	"les-randoms/utils"
	"strconv"
)

var summonerSpellsArray []SummonerSpell

type SummonerSpellsInfo struct { // Only 'useful' informations are parsed from JSON
	Type    string `json:"type"`
	Version string `json:"version"`
	Data    struct {
		SummonerBarrier                  SummonerSpell `json:"SummonerBarrier"`
		SummonerBoost                    SummonerSpell `json:"SummonerBoost"`
		SummonerDot                      SummonerSpell `json:"SummonerDot"`
		SummonerExhaust                  SummonerSpell `json:"SummonerExhaust"`
		SummonerFlash                    SummonerSpell `json:"SummonerFlash"`
		SummonerHaste                    SummonerSpell `json:"SummonerHaste"`
		SummonerHeal                     SummonerSpell `json:"SummonerHeal"`
		SummonerMana                     SummonerSpell `json:"SummonerMana"`
		SummonerPoroRecall               SummonerSpell `json:"SummonerPoroRecall"`
		SummonerPoroThrow                SummonerSpell `json:"SummonerPoroThrow"`
		SummonerSmite                    SummonerSpell `json:"SummonerSmite"`
		SummonerSnowURFSnowball_Mark     SummonerSpell `json:"SummonerSnowURFSnowball_Mark"`
		SummonerSnowball                 SummonerSpell `json:"SummonerSnowball"`
		SummonerTeleport                 SummonerSpell `json:"SummonerTeleport"`
		Summoner_UltBookPlaceholder      SummonerSpell `json:"Summoner_UltBookPlaceholder"`
		Summoner_UltBookSmitePlaceholder SummonerSpell `json:"Summoner_UltBookSmitePlaceholder"`
	} `json:"data"`
}

type SummonerSpell struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Image       struct {
		Full   string `json:"full"`
		Sprite string `json:"sprite"`
		Group  string `json:"group"`
		X      int    `json:"x"`
		Y      int    `json:"y"`
		W      int    `json:"w"`
		H      int    `json:"h"`
	} `json:"image"`
}

func GetAllSummonerSpells() (*SummonerSpellsInfo, error) {
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

func GetSummonerSpellImageNameByKey(key int) string {
	for _, summonerSpell := range getSummonerSpellsArray() {
		if summonerSpell.Key == strconv.Itoa(key) {
			return summonerSpell.Image.Full
		}
	}
	utils.LogError("Summoner Spell Key Not Found : " + strconv.Itoa(key))
	return ""
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

func getSummonerSpellsArray() []SummonerSpell {
	updateServerInfoIfNecessary()
	return summonerSpellsArray
}

func updateServerSummonerSpellsInfo() error {
	summonerSpellsInfo, err := GetAllSummonerSpells()
	if err != nil {
		utils.LogError("Error while refreshing server summoner spells info :\n" + err.Error())
		return err
	}
	summonerSpellsArray = make([]SummonerSpell, 0)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerBarrier)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerBoost)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerDot)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerFlash)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerHaste)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerHeal)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerMana)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerPoroRecall)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerPoroThrow)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerSmite)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerSnowURFSnowball_Mark)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerSnowball)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.SummonerTeleport)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.Summoner_UltBookPlaceholder)
	summonerSpellsArray = append(summonerSpellsArray, summonerSpellsInfo.Data.Summoner_UltBookSmitePlaceholder)
	return nil
}
