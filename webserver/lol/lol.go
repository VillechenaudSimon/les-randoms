package lol

import webserver "les-randoms/webserver/logic"

func SetupRoutes() {
	lol := webserver.Router.Group("/lol")

	aram := lol.Group("/aram")
	aram.GET("", handleAramRoute)
	aram.POST("", handleAramRoute)
	aram.GET("/:subNavItem", handleAramRoute)
	aram.POST("/:subNavItem", handleAramRoute)

	players := lol.Group("/players")
	players.GET("", handlePlayersRoute)
	players.POST("", handlePlayersRoute)
	players.GET("/:subNavItem", handlePlayersRoute)
	players.POST("/:subNavItem", handlePlayersRoute)
	players.GET("/:subNavItem/:param1", handlePlayersRoute)
	players.POST("/:subNavItem/:param1", handlePlayersRoute)
}

type lolGameReviewData struct {
	GameMode     string
	GameDuration string
	BlueTeam     lolTeamGameReviewData
	RedTeam      lolTeamGameReviewData
}

type lolTeamGameReviewData struct {
	Players []lolPlayerGameReviewData
}

type lolPlayerGameReviewData struct {
	Version         string
	SummonerName    string
	ChampionName    string
	SummonerSpell1  string
	SummonerSpell2  string
	CreepScore      int
	GoldEarned      string
	Kills           int
	Deaths          int
	Assists         int
	KDA             string
	WardsPlaced     int
	WardsKilled     int
	PinkWardsBought int
	VisionScore     int
	Trinket         string
	Items           []string
}

type lolProfileData struct {
	Version  string
	Summoner struct {
		Name        string
		IconId      int
		Level       int
		SoloDuoRank string
		SoloDuoLP   int
	}
}
