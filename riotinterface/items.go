package riotinterface

import (
	"encoding/json"
	"les-randoms/utils"
	"strconv"
)

var itemsMap map[string]Item

type ItemsInfo struct { // Only 'useful' informations are parsed from JSON
	Type    string          `json:"type"`
	Version string          `json:"version"`
	Data    map[string]Item `json:"data"`
}

type Item struct { // Only 'useful' informations are parsed from JSON
	Name      string    `json:"name"`
	PlainText string    `json:"plaintext"`
	Image     riotImage `json:"image"`
}

func GetItemsMap() map[string]Item {
	updateServerInfoIfNecessary()
	return itemsMap
}

func GetItemsFromInt(key int) Item {
	return GetItemsMap()[strconv.Itoa(key)]
}

func getItemsInfo() (*ItemsInfo, error) {
	body, err := getAllItemsJSON()
	if err != nil {
		return nil, err
	}
	itemsInfo, err := parseItemsJSON(body)
	if err != nil {
		return nil, err
	}
	return itemsInfo, nil
}

func getAllItemsJSON() ([]byte, error) {
	return requestRIOTAPI("https://ddragon.leagueoflegends.com/cdn/" + GetPatch() + "/data/en_US/item.json")
}

func parseItemsJSON(body []byte) (*ItemsInfo, error) {
	data := &ItemsInfo{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	return data, nil
}

func updateServerItemsInfo() error {
	itemsInfo, err := getItemsInfo()
	if err != nil {
		utils.LogError("Error while updating server items info :\n" + err.Error())
		return err
	}
	itemsMap = itemsInfo.Data
	return nil
}
