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
		SummonerMissingInDB: "No Summoner match the request",
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

func Summoner_SelectAll(queryPart string) ([]Summoner, error) {
	rows, err := SelectDatabase("summonerId, userId, accountId, puuid, name, profileIconId, level, revisionDate, lastUpdated FROM Summoner " + queryPart)
	defer rows.Close()
	if err != nil {
		utils.LogError("Error while selecting on Summoner table : " + err.Error())
		return nil, err
	}
	summoners := make([]Summoner, 0)
	for rows.Next() {
		var summonerId string
		var userId int
		var accountId string
		var puuid string
		var name string
		var profileIconId int
		var level int
		var revisionDate int64
		var lastUpdated []uint8
		err = rows.Scan(&summonerId, &userId, &accountId, &puuid, &name, &profileIconId, &level, &revisionDate, &lastUpdated)
		if err != nil {
			utils.LogError("Error while scanning on Summoner table : " + err.Error())
			return nil, err
		}
		parsedLastUpdated, err := time.Parse(utils.DBDateTimeFormat, string(lastUpdated))
		if err != nil {
			utils.LogError("Error while parsing a summoner date : " + err.Error())
			continue
		}
		summoners = append(summoners, Summoner{
			SummonerId:    summonerId,
			UserId:        userId,
			AccountId:     accountId,
			Puuid:         puuid,
			Name:          name,
			ProfileIconId: profileIconId,
			Level:         level,
			RevisionDate:  revisionDate,
			LastUpdated:   parsedLastUpdated,
		})
	}
	return summoners, nil
}

func Summoner_SelectFirst(queryPart string) (Summoner, error) {
	rows, err := SelectDatabase("summonerId, userId, accountId, puuid, name, profileIconId, level, revisionDate, lastUpdated FROM Summoner " + queryPart)
	if err != nil {
		utils.LogError("Error while selecting on Summoner table : " + err.Error())
		return Summoner{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return Summoner{}, errors.New(SummonerErrors.SummonerMissingInDB)
	}
	var summonerId string
	var userId int
	var accountId string
	var puuid string
	var name string
	var profileIconId int
	var level int
	var revisionDate int64
	var lastUpdated []uint8
	err = rows.Scan(&summonerId, &userId, &accountId, &puuid, &name, &profileIconId, &level, &revisionDate, &lastUpdated)
	if err != nil {
		utils.LogError("Error while scanning on Summoner table : " + err.Error())
		return Summoner{}, err
	}
	parsedLastUpdated, err := time.Parse(utils.DBDateTimeFormat, string(lastUpdated))
	if err != nil {
		utils.LogError("Error while parsing a summoner date : " + err.Error())
		return Summoner{}, err
	}
	return Summoner{
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
	result, err := InsertDatabase("Summoner(summonerId, userId, accountId, puuid, name, profileIconId, level, revisionDate, lastUpdated) VALUES(" +
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
	result, err := UpdateDatabase("Summoner SET " +
		"userId=" + utils.Esc(fmt.Sprint(summoner.UserId)) + ", " +
		"accountId=" + utils.Esc(summoner.AccountId) + ", " +
		"puuid=" + utils.Esc(summoner.Puuid) + ", " +
		"name=" + utils.Esc(summoner.Name) + ", " +
		"profileIconId=" + utils.Esc(fmt.Sprint(summoner.ProfileIconId)) + ", " +
		"level=" + utils.Esc(fmt.Sprint(summoner.Level)) + ", " +
		"revisionDate=" + utils.Esc(fmt.Sprint(summoner.RevisionDate)) + ", " +
		"lastUpdated=" + utils.Esc(time.Now().UTC().Format(utils.DBDateTimeFormat)) + " " +
		"FROM (SELECT * FROM Summoner WHERE summonerId=" + utils.Esc(summoner.SummonerId) + ")")
	return result, err
}
