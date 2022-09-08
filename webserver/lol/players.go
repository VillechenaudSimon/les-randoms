package lol

import (
	"fmt"
	radbwrapper "les-randoms/radb-wrapper"
	"les-randoms/riotinterface"
	"les-randoms/utils"
	webserver "les-randoms/webserver/logic"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handlePlayersRoute(c *gin.Context) {
	session := webserver.GetSession(c)

	if c.Param("subNavItem") == "" {
		c.Redirect(http.StatusFound, "/lol/players/Profile")
	}

	data := playersData{}

	webserver.SetupNavData(&data.LayoutData.NavData, session)

	selectedItemName := webserver.SetupSubnavData(&data.LayoutData.SubnavData, c, "Player Analyser", []string{"Profile(WIP)🚧", "LastGame", "Ladder", "LadderChampPool(WIP)🚧"}, map[string]string{"Profile(WIP)🚧": "Profile (WIP)🚧", "LastGame": "Last Game", "Ladder": "Ladder", "LadderChampPool(WIP)🚧": "Ladder Champ Pool (WIP)🚧"})

	webserver.SetupContentHeaderData(&data.ContentHeaderData, session)
	data.ContentHeaderData.Title = selectedItemName

	switch data.LayoutData.SubnavData.SelectedSubnavItemIndex {
	case 0:
		data.ProfileParameters.SummonerName = c.Param("param1")
		if setupLolProfileData(&data) != nil {
			c.Redirect(http.StatusFound, "/lol/players/Profile")
		}
	case 1:
		data.LastGameParameters.SummonerName = c.Param("param1")
		if setupLolLastGameData(&data) != nil {
			c.Redirect(http.StatusFound, "/lol/players/LastGame")
		}
	case 2:
		if setupLadderTableData(&data) != nil {
			c.Redirect(http.StatusFound, "/")
		}
	case 3:
		if setupLadderChampPoolTableData(&data) != nil {
			c.Redirect(http.StatusFound, "/lol/players/LadderChampPool")
		}
	}

	c.HTML(http.StatusFound, "players.tmpl.html", data)
}

type playersData struct {
	LayoutData        webserver.LayoutData
	ContentHeaderData webserver.ContentHeaderData
	ProfileParameters struct {
		SummonerName string
	}
	LastGameParameters struct {
		SummonerName string
	}
	LolGameReviewData        lolGameReviewData
	LolProfileData           lolProfileData
	LadderChampPoolTableData webserver.CustomTableData
	LadderTableData          webserver.CustomTableData
}

func setupLolProfileData(data *playersData) error {
	if data.ProfileParameters.SummonerName == "" {
		return nil
	}

	data.LolProfileData.Version = riotinterface.GetPatch()

	data.LolProfileData.Summoner.Name = data.ProfileParameters.SummonerName
	summoner, _, err := radbwrapper.GetSummonerFromName(data.LolProfileData.Summoner.Name)
	if err != nil {
		utils.LogError(err.Error())
		return err
	}
	data.LolProfileData.Summoner.Level = summoner.Level
	data.LolProfileData.Summoner.IconId = summoner.ProfileIconId
	leagueEntry, err := riotinterface.GetSoloDuoEntryBySummonerId(summoner.SummonerId)
	if err == nil {
		data.LolProfileData.Summoner.SoloDuo.TierFlat = leagueEntry.Tier
		data.LolProfileData.Summoner.SoloDuo.RankFlat = leagueEntry.Rank
		data.LolProfileData.Summoner.SoloDuo.TierRankDisplay = riotinterface.ParseTierRank(leagueEntry.Tier, leagueEntry.Rank)
		data.LolProfileData.Summoner.SoloDuo.LP = leagueEntry.LeaguePoints
	} else if err.Error() == riotinterface.LeagueV4Errors.MissingSummonerSoloDuoEntry {
		data.LolProfileData.Summoner.SoloDuo.TierRankDisplay = "Unranked"
	} else {
		utils.LogError(err.Error())
	}

	// TODO
	data.LolProfileData.Summoner.SoloDuo.Roles = append(data.LolProfileData.Summoner.SoloDuo.Roles, lolProfileDataSummonerRoles{"Top", 10, 5, 5})
	data.LolProfileData.Summoner.SoloDuo.Roles = append(data.LolProfileData.Summoner.SoloDuo.Roles, lolProfileDataSummonerRoles{"Jungle", 20, 9, 11})
	data.LolProfileData.Summoner.SoloDuo.Roles = append(data.LolProfileData.Summoner.SoloDuo.Roles, lolProfileDataSummonerRoles{"Mid", 30, 12, 18})
	data.LolProfileData.Summoner.SoloDuo.Roles = append(data.LolProfileData.Summoner.SoloDuo.Roles, lolProfileDataSummonerRoles{"Bot", 8, 8, 0})
	data.LolProfileData.Summoner.SoloDuo.Roles = append(data.LolProfileData.Summoner.SoloDuo.Roles, lolProfileDataSummonerRoles{"Support", 11, 6, 5})

	// TODO
	data.LolProfileData.Summoner.SoloDuo.GamesCount = 79
	data.LolProfileData.Summoner.SoloDuo.WinsCount = 40
	data.LolProfileData.Summoner.SoloDuo.LossesCount = 39

	matchIds, err := riotinterface.GetLastsMatchsIdsFromPuuid(summoner.Puuid, "5")
	if err != nil {
		utils.LogError(err.Error())
	} else {
		for _, mId := range matchIds {
			match, err := riotinterface.GetMatchFromId(fmt.Sprint(mId))
			if err != nil {
				utils.LogError(err.Error())
			} else {
				g := lolProfileGame{}
				g.Info.GameDuration = formatGameDuration(match.Info.GameDuration)
				g.Info.GameMode = riotinterface.ParseGameModeFromQueueId(match.Info.QueueId)

				for _, p := range match.Info.Participants {
					if p.Puuid == summoner.Puuid {
						g.Info.IsWin = p.Win
						g.Player.Champion = p.ChampionName
						g.Player.Summoners = []string{riotinterface.GetSummonerSpellImageNameByKey(p.Summoner1Id), riotinterface.GetSummonerSpellImageNameByKey(p.Summoner2Id)}
						g.Player.Runes = make([]string, 6)
						/*for _, r := range p.Perks.Styles {

						}*/
						break
					}
				}
				data.LolProfileData.Games = append(data.LolProfileData.Games, g)
			}
		}
	}

	return nil
}

func setupLolLastGameData(data *playersData) error {
	if data.LastGameParameters.SummonerName == "" {
		return nil
	}

	puuid, err := riotinterface.GetPuuidFromSummonerName(data.LastGameParameters.SummonerName)
	if err != nil {
		utils.LogError(err.Error())
		return err
	}
	matchId, err := riotinterface.GetLastMatchIdFromPuuid(puuid)
	if err != nil {
		utils.LogError(err.Error())
		return err
	}
	match, err := riotinterface.GetMatchFromId(matchId)
	if err != nil {
		utils.LogError(err.Error())
		return err
	}
	data.LolGameReviewData.GameDuration = formatGameDuration(match.Info.GameDuration)
	data.LolGameReviewData.GameMode = riotinterface.ParseGameModeFromQueueId(match.Info.QueueId)

	for _, p := range match.Info.Participants {
		var kda string
		if p.Deaths > 0 {
			tmp := (float32(p.Kills) + float32(p.Assists)) / float32(p.Deaths)
			kda = strconv.Itoa(int(tmp)) + "." + strconv.Itoa(int(tmp*100)%100)
		} else {
			kda = "Perfect"
		}
		trinket := riotinterface.GetItemsFromInt(p.Item6).Image.Full
		items := make([]string, 0)
		items = append(items, riotinterface.GetItemsFromInt(p.Item0).Image.Full)
		items = append(items, riotinterface.GetItemsFromInt(p.Item1).Image.Full)
		items = append(items, riotinterface.GetItemsFromInt(p.Item2).Image.Full)
		items = append(items, riotinterface.GetItemsFromInt(p.Item3).Image.Full)
		items = append(items, riotinterface.GetItemsFromInt(p.Item4).Image.Full)
		items = append(items, riotinterface.GetItemsFromInt(p.Item5).Image.Full)
		player := lolPlayerGameReviewData{
			"",
			p.SummonerName,
			p.ChampionName,
			riotinterface.GetSummonerSpellImageNameByKey(p.Summoner1Id),
			riotinterface.GetSummonerSpellImageNameByKey(p.Summoner2Id),
			p.TotalMinionsKilled,
			strconv.Itoa(p.GoldEarned/1000) + "." + strconv.Itoa((p.GoldEarned%1000)/100),
			p.Kills,
			p.Deaths,
			p.Assists,
			kda,
			p.WardsPlaced,
			p.WardsKilled,
			p.VisionWardsBoughtInGame,
			p.VisionScore,
			trinket,
			items,
		}
		player.Version, _ = riotinterface.GetLastVersionFromGameVersion(match.Info.GameVersion)
		if p.TeamId == 100 {
			data.LolGameReviewData.BlueTeam.Players = append(data.LolGameReviewData.BlueTeam.Players, player)
		} else { // p.TeamId == 200
			data.LolGameReviewData.RedTeam.Players = append(data.LolGameReviewData.RedTeam.Players, player)
		}
	}

	return nil
}

func setupLadderTableData(data *playersData) error {
	data.LadderTableData.HeaderList = []string{"Summoner Icon", "LP", "Summoner Name"}
	data.LadderTableData.ColumnTypes = []webserver.CustomTableColumnType{webserver.CustomTableColumnTypeImage, webserver.CustomTableColumnTypeNumber, webserver.CustomTableColumnTypeText}
	challengerLeague, err := riotinterface.GetSoloDuoChallengerLeague()
	if err != nil {
		utils.LogError(err.Error())
		return err
	}
	summoners, err := radbwrapper.LeagueListToSummoners(challengerLeague)
	if err != nil {
		utils.LogError(err.Error())
		return err
	}
	for _, entry := range challengerLeague.Entries {
		data.LadderTableData.ItemList = append(
			data.LadderTableData.ItemList,
			webserver.TableItemData{FieldList: []string{
				riotinterface.GetProfileIconUrl(summoners[entry.SummonerId].ProfileIconId),
				fmt.Sprint(entry.LeaguePoints),
				entry.SummonerName,
			}},
		)
	}
	data.LadderTableData.SortColumnIndex = 1
	data.LadderTableData.SortOrder = 0

	return nil
}

func setupLadderChampPoolTableData(data *playersData) error {
	data.LadderChampPoolTableData.HeaderList = append(data.LadderChampPoolTableData.HeaderList, "Work In Progress..")
	data.LadderChampPoolTableData.ColumnTypes = []webserver.CustomTableColumnType{webserver.CustomTableColumnTypeText}
	return nil
}

func formatGameDuration(gameDuration int64) string {
	minutes := strconv.Itoa(int(gameDuration) / 60)
	if int(gameDuration)/60 < 10 {
		minutes = "0" + minutes
	}
	seconds := strconv.Itoa(int(gameDuration) % 60)
	if int(gameDuration)%60 < 10 {
		seconds = "0" + seconds
	}
	return minutes + ":" + seconds
}
