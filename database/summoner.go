package database

import (
	"database/sql"
	"errors"
	"fmt"
	"les-randoms/utils"
	"reflect"
	"time"
)

var SummonerErrors SummonerErrorsConst

func init() {
	SummonerErrors = SummonerErrorsConst{
		SummonerMissingInDB: "No " + databaseTableNames.Summoner + " match the request",
	}
}

type SummonerErrorsConst struct {
	SummonerMissingInDB string
}

type Summoner struct {
	SummonerId    string
	UserId        int
	AccountId     string
	Puuid         string
	Name          string
	ProfileIconId int
	Level         int
	RevisionDate  int64
	LastUpdated   time.Time
}

func Summoner_Construct(summonerId string, userId int, accountId string, puuid string, name string, profileIconId int, level int, revisionDate int64) Summoner {
	return Summoner{
		SummonerId:    summonerId,
		UserId:        userId,
		AccountId:     accountId,
		Puuid:         puuid,
		Name:          name,
		ProfileIconId: profileIconId,
		Level:         level,
		RevisionDate:  revisionDate,
		LastUpdated:   time.Now().UTC(),
	}
}

func (summoner Summoner) ToStringSlice() []string {
	return []string{
		summoner.SummonerId,
		fmt.Sprint(summoner.UserId),
		summoner.AccountId,
		summoner.Puuid,
		summoner.Name,
		fmt.Sprint(summoner.ProfileIconId),
		fmt.Sprint(summoner.Level),
		fmt.Sprint(summoner.RevisionDate),
		summoner.LastUpdated.Format(utils.DateTimeFormat),
	}
}

func Summoners_ToDBStructSlice(summoners []Summoner) []DBStruct {
	var r []DBStruct
	for _, summoner := range summoners {
		r = append(r, summoner)
	}
	return r
}

func Summoner_GetType() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(&Summoner{})).Type()
}

func summoner_SelectQuery(queryPart string) (*sql.Rows, error) {
	rows, err := SelectDatabase("summonerid, userid, accountid, puuid, name, profileiconid, level, revisiondate, lastupdated FROM " + databaseTableNames.Summoner + " " + queryPart)
	if err != nil {
		utils.LogError("Error while selecting on " + databaseTableNames.Summoner + " table : " + err.Error())
		return nil, err
	}
	return rows, nil
}

func summoner_ScanOne(rows *sql.Rows) (bool, Summoner, error) {
	if rows.Next() {
		var summonerId string
		var userId int
		var accountId string
		var puuid string
		var name string
		var profileIconId int
		var level int
		var revisionDate int64
		var lastUpdated []uint8
		err := rows.Scan(&summonerId, &userId, &accountId, &puuid, &name, &profileIconId, &level, &revisionDate, &lastUpdated)
		if err != nil {
			utils.LogError("Error while scanning on " + databaseTableNames.Summoner + " table : " + err.Error())
			return true, Summoner{}, err
		}
		parsedLastUpdated, err := time.Parse(utils.DBDateTimeFormat, string(lastUpdated))
		if err != nil {
			utils.LogError("Error while parsing a " + databaseTableNames.Summoner + " date : " + err.Error())
		}
		return true, Summoner{
			SummonerId:    summonerId,
			UserId:        userId,
			AccountId:     accountId,
			Puuid:         puuid,
			Name:          name,
			ProfileIconId: profileIconId,
			Level:         level,
			RevisionDate:  revisionDate,
			LastUpdated:   parsedLastUpdated,
		}, nil
	} else {
		return false, Summoner{}, nil
	}
}

func Summoner_SelectAll(queryPart string) ([]Summoner, error) {
	rows, err := summoner_SelectQuery(queryPart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	summoners := make([]Summoner, 0)
	b, s, err := summoner_ScanOne(rows)
	for b {
		if err == nil {
			summoners = append(summoners, s)
		}
		b, s, err = summoner_ScanOne(rows)
	}
	return summoners, nil
}

func Summoner_SelectAllInMapId(queryPart string) (map[string]Summoner, error) {
	rows, err := summoner_SelectQuery(queryPart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	summoners := make(map[string]Summoner, 0)
	b, s, err := summoner_ScanOne(rows)
	for b {
		if err == nil {
			summoners[s.SummonerId] = s
		}
		b, s, err = summoner_ScanOne(rows)
	}
	return summoners, nil
}

func Summoner_SelectAllInMapName(queryPart string) (map[string]Summoner, error) {
	rows, err := summoner_SelectQuery(queryPart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	summoners := make(map[string]Summoner, 0)
	b, s, err := summoner_ScanOne(rows)
	for b {
		if err == nil {
			summoners[s.Name] = s
		}
		b, s, err = summoner_ScanOne(rows)
	}
	return summoners, nil
}

func Summoner_SelectFirst(queryPart string) (Summoner, error) {
	rows, err := summoner_SelectQuery(queryPart)
	if err != nil {
		return Summoner{}, err
	}
	defer rows.Close()
	b, s, err := summoner_ScanOne(rows)
	if err != nil {
		return Summoner{}, err
	}
	if b {
		return s, nil
	}
	return Summoner{}, errors.New(SummonerErrors.SummonerMissingInDB)
}

func Summoner_CreateNew(summonerId string, userId int, accountId string, puuid string, name string, profileIconId int, level int, revisionDate int64) (Summoner, sql.Result, error) {
	summoner := Summoner{
		SummonerId:    summonerId,
		UserId:        userId,
		AccountId:     accountId,
		Puuid:         puuid,
		Name:          name,
		ProfileIconId: profileIconId,
		Level:         level,
		RevisionDate:  revisionDate,
	}
	result, err := Summoner_Insert(summoner)
	return summoner, result, err
}

func Summoner_Insert(summoner Summoner) (sql.Result, error) {
	result, err := InsertDatabase(databaseTableNames.Summoner + "(summonerid, userid, accountid, puuid, name, profileiconid, level, revisiondate, lastupdated) VALUES(" +
		utils.Esc(summoner.SummonerId) + ", " +
		utils.Esc(fmt.Sprint(summoner.UserId)) + ", " +
		utils.Esc(summoner.AccountId) + ", " +
		utils.Esc(summoner.Puuid) + ", " +
		utils.Esc(summoner.Name) + ", " +
		utils.Esc(fmt.Sprint(summoner.ProfileIconId)) + ", " +
		utils.Esc(fmt.Sprint(summoner.Level)) + ", " +
		utils.Esc(fmt.Sprint(summoner.RevisionDate)) + ", " +
		utils.Esc(time.Now().UTC().Format(utils.DBDateTimeFormat)) + ")")
	return result, err
}

func Summoner_Update(summoner Summoner) (sql.Result, error) {
	result, err := UpdateDatabase(databaseTableNames.Summoner + " SET " +
		"userid=" + utils.Esc(fmt.Sprint(summoner.UserId)) + ", " +
		"accountid=" + utils.Esc(summoner.AccountId) + ", " +
		"puuid=" + utils.Esc(summoner.Puuid) + ", " +
		"name=" + utils.Esc(summoner.Name) + ", " +
		"profileiconid=" + utils.Esc(fmt.Sprint(summoner.ProfileIconId)) + ", " +
		"level=" + utils.Esc(fmt.Sprint(summoner.Level)) + ", " +
		"revisiondate=" + utils.Esc(fmt.Sprint(summoner.RevisionDate)) + ", " +
		"lastupdated=" + utils.Esc(time.Now().UTC().Format(utils.DBDateTimeFormat)) + " " +
		"FROM (SELECT * FROM " + databaseTableNames.Summoner + " WHERE summonerid=" + utils.Esc(summoner.SummonerId) + ") as summonerstoupdate")
	return result, err
}
