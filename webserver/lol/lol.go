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
		Name    string
		IconId  int
		Level   int
		SoloDuo lolProfileDataQueueInfo
	}
	Games []lolProfileGame
}

type lolProfileGame struct {
	Info   lolProfileGameInfo
	Player lolProfileGamePlayer
	Teams  []lolProfileGameTeam
}

type lolProfileGameInfo struct {
	IsWin        bool
	GameMode     string
	GameDuration string
}

type lolProfileGamePlayer struct {
	Champion  string
	Summoners []string
	Runes     []lolProfileGamePlayerRune
	Build     []lolProfileGamePlayerSlot
}

type lolProfileGamePlayerRune struct {
	Name     string
	IconPath string
}

type lolProfileGamePlayerSlot struct {
	Used     bool
	Name     string
	IconPath string
}

type lolProfileGameTeam struct {
	Players []struct {
		Champion string
		Name     string
	}
}

type lolProfileDataQueueInfo struct {
	TierRankDisplay string // GrandMaster
	TierFlat        string // Bronze
	RankFlat        string // 2
	LP              int
	GamesCount      int
	WinsCount       int
	LossesCount     int
	Roles           []lolProfileDataSummonerRoles
}

type lolProfileDataSummonerRoles struct {
	Name        string
	GamesCount  int
	WinsCount   int
	LossesCount int
}
