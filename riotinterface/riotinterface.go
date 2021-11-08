package riotinterface

import (
	"les-randoms/utils"
	"net/http"
	"os"
)

func Test() {
	body, err := GetMatchInfo("EUW1_5540830822")
	if err != nil {
		utils.LogError(err.Error())
	}
	match, err := ParseMatchJSON(body)
	if err != nil {
		utils.LogError(err.Error())
	}
	utils.LogDebug(match.Metadata.MatchId)
}

func prepareRequestHeader(request *http.Request) {
	request.Header.Set("Accept-Charset", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("X-Riot-Token", os.Getenv("X_RIOT_TOKEN"))
}

func isResponseStatusNotOK(status string) bool {
	return status[0:3] != "200"
}
