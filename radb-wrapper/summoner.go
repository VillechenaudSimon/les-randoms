package radbwrapper

import (
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

func addRiotSummonerToDB(name string) (database.Summoner, error) {
	riotSummoner, err := riotinterface.GetSummonerFromName(name)
	if err != nil {
		return database.Summoner{}, err
	}
	summoner := riotSummonerToDBSummoner(*riotSummoner)
	_, err = database.Summoner_Insert(summoner)
	utils.LogNotNilError(err)
	return summoner, nil
}

func updateRiotSummonerToDB(name string) (database.Summoner, error) {
	riotSummoner, err := riotinterface.GetSummonerFromName(name)
	if err != nil {
		return database.Summoner{}, err
	}
	summoner := riotSummonerToDBSummoner(*riotSummoner)
	_, err = database.Summoner_Update(summoner)
	utils.LogNotNilError(err)
	return summoner, err
}

func updateSummonerIfNeeded(summoner database.Summoner) (database.Summoner, error) {
	if time.Since(summoner.LastUpdated).Hours() > 8 {
		summonerFromRiot, err := updateRiotSummonerToDB(summoner.Name)
		if err != nil { // In this case we return the last informations we have in the DB even if they are not the most recents possible
			return summoner, err
		}
		return summonerFromRiot, nil
	}
	return summoner, nil
}

func GetSummonerFromName(name string) (database.Summoner, error) {
	summoner, err := database.Summoner_SelectFirst("WHERE name=" + utils.Esc(name))
	if err == nil {
		/*
			if time.Since(summoner.LastUpdated).Hours() > 8 {
				summonerFromRiot, err := updateRiotSummonerToDB(name)
				if err != nil { // In this case we return the last informations we have in the DB even if they are not the most recents possible
					return summoner, err
				}
				summoner = summonerFromRiot
			}*/
		summoner, err = updateSummonerIfNeeded(summoner)
		utils.LogNotNilError(err)
		return summoner, err
	} else if err.Error() == database.SummonerErrors.SummonerMissingInDB {
		summoner, err = addRiotSummonerToDB(name)
		if err != nil {
			return database.Summoner{}, err
		}
	} else {
		return database.Summoner{}, err
	}
	return summoner, nil
}

func GetSummonersFromNames(namesGetter func(int) (bool, string)) (map[string]database.Summoner, error) {
	missingSumData := make(map[string]bool, 0)
	queryBody := "WHERE"
	s := ""
	i := 0
	b := true
	for {
		b, s = namesGetter(i)
		i++
		if !b {
			break
		}
		queryBody += " name=" + utils.Esc(s) + " OR"
		missingSumData[s] = true
	}
	summoners, err := database.Summoner_SelectAllInMapName(queryBody[:len(queryBody)-3])
	if err != nil {
		return nil, err
	}

	for n := range summoners {
		missingSumData[n] = false
	}

	for n, b := range missingSumData {
		if b {
			utils.LogDebug("AddSumDb")
			summoners[n], err = addRiotSummonerToDB(n)
		} else {
			utils.LogDebug("UpdateDb")
			summoners[n], err = updateSummonerIfNeeded(summoners[n])
		}
		utils.LogNotNilError(err)
	}

	return summoners, nil
}
