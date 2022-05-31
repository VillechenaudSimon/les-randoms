package radbwrapper

import (
	"les-randoms/database"
	"les-randoms/riotinterface"
	"les-randoms/utils"
	"time"
)

func riotSummonerToDBSummoner(summoner *riotinterface.Summoner) database.Summoner {
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

func addRiotSummonerToDBFromId(id string) (database.Summoner, error) {
	riotSummoner, err := riotinterface.GetSummonerFromId(id)
	if err != nil {
		return database.Summoner{}, err
	}
	summoner := riotSummonerToDBSummoner(riotSummoner)
	_, err = database.Summoner_Insert(summoner)
	utils.LogNotNilError(err)
	return summoner, nil
}

func addRiotSummonerToDBFromName(name string) (database.Summoner, error) {
	riotSummoner, err := riotinterface.GetSummonerFromName(name)
	if err != nil {
		return database.Summoner{}, err
	}
	summoner := riotSummonerToDBSummoner(riotSummoner)
	_, err = database.Summoner_Insert(summoner)
	utils.LogNotNilError(err)
	return summoner, nil
}

func updateRiotSummonerToDBFromId(id string) (database.Summoner, error) {
	riotSummoner, err := riotinterface.GetSummonerFromId(id)
	if err != nil {
		return database.Summoner{}, err
	}
	summoner := riotSummonerToDBSummoner(riotSummoner)
	_, err = database.Summoner_Update(summoner)
	if err != nil {
		return database.Summoner{}, err
	} else {
		return summoner, nil
	}
}

/*
func updateRiotSummonerToDBFromName(name string) (database.Summoner, error) {
	riotSummoner, err := riotinterface.GetSummonerFromName(name)
	if err != nil {
		return database.Summoner{}, err
	}
	summoner := riotSummonerToDBSummoner(*riotSummoner)
	_, err = database.Summoner_Update(summoner)
	utils.LogNotNilError(err)
	return summoner, err
}
*/

func updateSummonerIfNeeded(summoner database.Summoner) (database.Summoner, bool, error) {
	if time.Now().UTC().Sub(summoner.LastUpdated) > LadderSummonersUpdateSpacing {
		summonerFromRiot, err := updateRiotSummonerToDBFromId(summoner.SummonerId)
		if err != nil { // In this case we return the last informations we have in the DB even if they are not the most recents possible
			return summoner, false, err
		}
		return summonerFromRiot, true, nil
	}
	return summoner, false, nil
}

func GetSummonerFromName(name string) (database.Summoner, bool, error) {
	summoner, err := database.Summoner_SelectFirst("WHERE name=" + utils.Esc(name))
	if err == nil {
		summoner, b, err := updateSummonerIfNeeded(summoner)
		utils.LogNotNilError(err)
		return summoner, b, err
	} else if err.Error() == database.SummonerErrors.SummonerMissingInDB {
		summoner, err = addRiotSummonerToDBFromName(name)
		if err != nil {
			return database.Summoner{}, true, err
		}
		return summoner, true, nil
	} else {
		return database.Summoner{}, false, err
	}
}

/*
func GetSummonersFromNames(namesGetter func(int) (bool, string), riotAccess bool) (map[string]database.Summoner, error) {
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
			summoners[n], err = addRiotSummonerToDB(n)
		} else {
			summoners[n], _, err = updateSummonerIfNeeded(summoners[n])
		}
		utils.LogNotNilError(err)
	}

	return summoners, nil
}
*/

func GetSummonerFromId(id string) (database.Summoner, bool, error) {
	summoner, err := database.Summoner_SelectFirst("WHERE summonerid=" + utils.Esc(id))
	if err == nil {
		summoner, b, err := updateSummonerIfNeeded(summoner)
		return summoner, b, err
	} else if err.Error() == database.SummonerErrors.SummonerMissingInDB {
		summoner, err = addRiotSummonerToDBFromId(id)
		if err != nil {
			return database.Summoner{}, true, err
		}
		return summoner, true, nil
	} else {
		return database.Summoner{}, false, err
	}
}

func GetSummonersFromIds(idsGetter func(int) (bool, string), riotAccess bool) (map[string]database.Summoner, error) {
	missingSumData := make(map[string]bool, 0)
	queryBody := "WHERE"
	s := ""
	i := 0
	b := true
	for {
		b, s = idsGetter(i)
		i++
		if !b {
			break
		}
		queryBody += " summonerid=" + utils.Esc(s) + " OR"
		missingSumData[s] = true
	}
	summoners, err := database.Summoner_SelectAllInMapId(queryBody[:len(queryBody)-3])
	if err != nil {
		return nil, err
	}

	for n := range summoners {
		missingSumData[n] = false
	}

	for n, b := range missingSumData {
		var s database.Summoner
		if b {
			s, err = addRiotSummonerToDBFromId(n)
		} else {
			s, _, err = updateSummonerIfNeeded(summoners[n])
		}
		if err != nil {
			utils.LogError(err.Error())
		} else {
			summoners[n] = s
		}
	}
	return summoners, nil
}

func LeagueListToSummoners(league *riotinterface.LeagueListDTO) (map[string]database.Summoner, error) {
	return GetSummonersFromIds(func(i int) (bool, string) {
		if i < len(league.Entries) {
			return true, league.Entries[i].SummonerId
		}
		return false, ""
	}, true)
}
