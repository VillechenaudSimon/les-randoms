package riotinterface

import (
	"encoding/json"
	"les-randoms/utils"
)

func GetVersionsArray() ([]string, error) {
	body, err := getVersionsJSON()
	if err != nil {
		return nil, err
	}
	versions, err := parseVersionsJSON(body)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func GetLastVersion() (string, error) {
	versions, err := GetVersionsArray()
	if err != nil || len(versions) == 0 {
		return "ERROR", nil
	}
	return versions[0], nil
}

func getVersionsJSON() ([]byte, error) {
	return requestRIOTAPI("https://ddragon.leagueoflegends.com/api/versions.json")
}

func parseVersionsJSON(body []byte) ([]string, error) {
	data := make([]string, 0)
	err := json.Unmarshal(body, &data)
	if err != nil {
		utils.LogError(err.Error())
		return nil, err
	}
	return data, nil
}
