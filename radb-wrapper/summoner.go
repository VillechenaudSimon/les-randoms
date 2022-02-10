package radbwrapper

import (
	"fmt"
	"les-randoms/database"
	"les-randoms/riotinterface"
	"les-randoms/utils"
	"time"
)

func riotSummonerToDBSummoner(summoner riotinterface.Summoner) database.Summoner {
	return database.Summoner_Construct(
		summoner.Id,
		0,
		summoner.AccountId,
		summoner.Puuid,
		summoner.Name,
		summoner.ProfileIconId,
		summoner.SummonerLevel,
		summoner.RevisionDate,
	)
}

func GetSummonerFromName(name string) (database.Summoner, error) {
	summoner, err := database.Summoner_SelectFirst("WHERE name=" + utils.Esc(name))
	if err == nil {
		utils.LogDebug("NOW : " + time.Now().Format(utils.DateTimeFormat) + " - " + time.Now().Location().String())
		utils.LogDebug("DB : " + summoner.LastUpdated.Format(utils.DateTimeFormat) + " - " + summoner.LastUpdated.Location().String())
		utils.LogDebug("HOURS : " + fmt.Sprint(time.Now().Sub(summoner.LastUpdated).Hours()))
		if time.Now().Sub(summoner.LastUpdated).Hours() > 1 {
			riotSummoner, err := riotinterface.GetSummonerFromName(name)
			if err != nil {
				return database.Summoner{}, err
			}
			summoner = riotSummonerToDBSummoner(*riotSummoner)
			database.Summoner_Update(summoner)
		}
	} else if err.Error() == database.SummonerErrors.SummonerMissingInDB {
		riotSummoner, err := riotinterface.GetSummonerFromName(name)
		if err != nil {
			return database.Summoner{}, err
		}
		summoner = riotSummonerToDBSummoner(*riotSummoner)
		_, err = database.Summoner_Insert(summoner)
		if err != nil {
			return database.Summoner{}, err
		}
	} else {
		return database.Summoner{}, err
	}
	return summoner, nil
}
