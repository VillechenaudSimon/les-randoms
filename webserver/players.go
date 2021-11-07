package webserver

import (
	"les-randoms/riotinterface"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlePlayersRoute(c *gin.Context) {
	session := getSession(c)

	data := playersData{}

	setupNavData(&data.LayoutData.NavData, session)

	selectedItemName := setupSubnavData(&data.LayoutData.SubnavData, c, "Player Analyser", []string{"Last Game"})

	setupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = selectedItemName

	data.LolGameReviewData.BlueTeam.Players = append(data.LolGameReviewData.BlueTeam.Players, lolPlayerGameReviewData{""})

	riotinterface.Test()

	c.HTML(http.StatusOK, "players.tmpl.html", data)
}
