package database

import (
	"database/sql"
	"errors"
	"fmt"
	"les-randoms/utils"
	"reflect"
	"time"
)

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
	rows, err := SelectDatabase("summonerId, userId, accountId, puuid, name, profileIconId, level, revisionDate FROM Summoner " + queryPart)
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
		return Summoner{}, errors.New("No Summoner match the request")
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

func Summoner_CreateNew(summonerId string, userId int, accountId string, puuid string, name string, profileIconId int, level int, revisionDate int64, lastUpdated time.Time) (Summoner, sql.Result, error) {
	result, err := InsertDatabase("Summoner(summonerId, userId, accountId, puuid, name, profileIconId, level, revisionDate, lastUpdated) VALUES(" +
		utils.Esc(summonerId) + ", " +
		utils.Esc(fmt.Sprint(userId)) + ", " +
		utils.Esc(accountId) + ", " +
		utils.Esc(puuid) + ", " +
		utils.Esc(name) + ", " +
		utils.Esc(fmt.Sprint(profileIconId)) + ", " +
		utils.Esc(fmt.Sprint(level)) + ", " +
		utils.Esc(fmt.Sprint(revisionDate)) +
		utils.Esc(lastUpdated.Format(utils.DBDateTimeFormat)) + ")")
	if err != nil {
		return Summoner{}, result, err
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
		LastUpdated:   lastUpdated,
	}, result, err
}
