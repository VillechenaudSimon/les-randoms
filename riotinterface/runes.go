package riotinterface

import (
	"encoding/json"
	"les-randoms/utils"
	"strings"
)

var runesMap map[int]Rune

type Rune struct { // Only 'useful' informations are parsed from JSON
	Id                      int    `json:"id"`
	Name                    string `json:"name"`
	MajorChangePatchVersion string `json:"majorChangePatchVersion"`
	IconPath                string `json:"iconPath"`
}

func GetRunesMap() map[int]Rune {
	updateServerInfoIfNecessary()
	return runesMap
}

func getRunesInfo() (*map[int]Rune, error) {
	body, err := getAllRunesJSON()
	if err != nil {
		return nil, err
	}
	runes, err := parseRunesJSON(body)
	if err != nil {
		return nil, err
	}
	return runes, nil
}

func getAllRunesJSON() ([]byte, error) {
	return requestRIOTAPI("https://raw.communitydragon.org/" + GetPatchShort() + "/plugins/rcp-be-lol-game-data/global/default/v1/perks.json")
}

func parseRunesJSON(body []byte) (*map[int]Rune, error) {
	data := &[]Rune{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	formatted := make(map[int]Rune, len(*data))
	for _, d := range *data {
		d.IconPath = strings.Join(strings.Split(d.IconPath, "/")[6:], "/")
		formatted[d.Id] = d
	}
	return &formatted, nil
}

func updateServerRunesInfo() error {
	newRunesInfo, err := getRunesInfo()
	if err != nil {
		utils.LogError("Error while updating server runes info :\n" + err.Error())
		return err
	}
	runesMap = *newRunesInfo
	return nil
}
