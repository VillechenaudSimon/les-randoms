package riotinterface

import (
	"encoding/json"
	"errors"
	"les-randoms/utils"
	"strings"
)

func GetVersionsCache() ([]string, error) {
	err := updateServerInfoIfNecessary()
	return VersionsList, err
}

func GetVersionsArray() ([]string, error) {
	body, err := getVersionsJSON()
	if err != nil {
		return nil, err
	}
	versions, err := parseVersionsJSON(body)
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, errors.New("versions array is empty")
	}
	return versions, nil
}

func GetLastVersion() (string, error) {
	versions, err := GetVersionsCache()
	if err != nil {
		return "ERROR", err
	}
	return versions[0], nil
}

func GetLastVersionFromGameVersion(gameVersion string) (string, error) {
	versions, err := GetVersionsCache()
	if err != nil {
		return "ERROR", err
	}
	splits := strings.Split(gameVersion, ".")
	gameVersionStart := splits[0] + "." + splits[1]
	for _, v := range versions {
		splits = strings.Split(v, ".")
		if gameVersionStart == (splits[0] + "." + splits[1]) {
			return v, nil
		}
	}
	return "ERROR", errors.New("Game Version \"" + gameVersion + "\" not found in versions array")
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
