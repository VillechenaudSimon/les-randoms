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
			trinket := riotinterface.GetItemsFromInt(p.Item6).Image.Full
			items := make([]string, 0)
			items = append(items, riotinterface.GetItemsFromInt(p.Item0).Image.Full)
			items = append(items, riotinterface.GetItemsFromInt(p.Item1).Image.Full)
			items = append(items, riotinterface.GetItemsFromInt(p.Item2).Image.Full)
			items = append(items, riotinterface.GetItemsFromInt(p.Item3).Image.Full)
			items = append(items, riotinterface.GetItemsFromInt(p.Item4).Image.Full)
			items = append(items, riotinterface.GetItemsFromInt(p.Item5).Image.Full)
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
				trinket,
				items,
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
