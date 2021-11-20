package webserver

import (
	"les-randoms/riotinterface"
	"les-randoms/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handlePlayersRoute(c *gin.Context) {
	session := getSession(c)

	data := playersData{}

	setupNavData(&data.LayoutData.NavData, session)

	selectedItemName := setupSubnavData(&data.LayoutData.SubnavData, c, "Player Analyser", []string{"Last Game"})

	setupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = selectedItemName

	match, err := riotinterface.GetMatchFromId("EUW1_5540830822")
	if err != nil {
		utils.LogError(err.Error())
	} else {
		for _, p := range match.Info.Participants {
			kda := (float32(p.Kills) + float32(p.Assists)) / float32(p.Deaths)
			player := lolPlayerGameReviewData{
				p.ChampionName,
				riotinterface.GetSummonerSpellImageNameByKey(p.Summoner1Id),
				riotinterface.GetSummonerSpellImageNameByKey(p.Summoner2Id),
				p.TotalMinionsKilled,
				strconv.Itoa(p.GoldEarned/1000) + "." + strconv.Itoa((p.GoldEarned%1000)/100),
				p.Kills,
				p.Deaths,
				p.Assists,
				strconv.Itoa(int(kda)) + "." + strconv.Itoa(int(kda*100)%100),
				p.WardsPlaced,
				p.WardsKilled,
				p.VisionWardsBoughtInGame,
				p.VisionScore,
			}
			if p.TeamId == 100 {
				data.LolGameReviewData.BlueTeam.Players = append(data.LolGameReviewData.BlueTeam.Players, player)
			} else { // p.TeamId == 200
				data.LolGameReviewData.RedTeam.Players = append(data.LolGameReviewData.RedTeam.Players, player)
			}
		}
	}

	c.HTML(http.StatusOK, "players.tmpl.html", data)
}
