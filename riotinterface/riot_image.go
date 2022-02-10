package riotinterface

import "fmt"

type riotImage struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
}

func GetProfileIconUrl(profileIconId int) string {
	return "https://ddragon.leagueoflegends.com/cdn/" + GetPatch() + "/img/profileicon/" + fmt.Sprint(profileIconId) + ".png"
}
