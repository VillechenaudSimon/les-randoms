package riotinterface

var lastUpdateVersion string

type ItemsInfo struct { // Only 'useful' informations are parsed from JSON
	Type    string `json:"type"`
	Version string `json:"version"`
	Data    map[string]struct {
		Name string `json:"name"`
	} `json:"data"`
}

func updateServerItemsInfo() error {
	return nil
}
