package database

import (
	"encoding/json"
	"fmt"
	"les-randoms/riotinterface"
	"les-randoms/utils"
	"reflect"
	"time"
)

type LolMatch struct {
	Id            int
	MatchId       string
	Server        string
	StartDate     time.Time
	EndDate       time.Time
	JsonData      []byte
	dataProcessed bool
	data          riotinterface.Match
}

func (lm LolMatch) ToStringSlice() []string {
	return []string{fmt.Sprint(lm.Id), lm.MatchId, lm.Server, lm.StartDate.Format(utils.DBDateTimeFormat), lm.EndDate.Format(utils.DateTimeFormat), "[...]"}
}

func LolMatchs_ToDBStructSlice(lms []LolMatch) []DBStruct {
	var r []DBStruct
	for _, lm := range lms {
		r = append(r, lm)
	}
	return r
}

func LolMatch_GetType() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(&LolMatch{})).Type()
}

func LolMatch_SelectAll(queryPart string) ([]LolMatch, error) {
	rows, err := SelectDatabase("id, matchid, server, startdate, enddate FROM " + databaseTableNames.LolMatch + " " + queryPart)
	if err != nil {
		utils.LogError("Error while selecting on " + databaseTableNames.LolMatch + " table : " + err.Error())
		return nil, err
	}
	defer rows.Close()
	lms := make([]LolMatch, 0)
	for rows.Next() {
		var id int
		var matchid string
		var server string
		var startdate []uint8
		var enddate []uint8
		err = rows.Scan(&id, &matchid, &server, &startdate, &enddate)
		if err != nil {
			utils.LogError("Error while scanning a " + databaseTableNames.LolMatch + " : " + err.Error())
			continue
		}
		parsedStartDate, err := time.Parse(utils.DBDateTimeFormat, string(startdate))
		if err != nil {
			utils.LogError("Error while parsing a " + databaseTableNames.LolMatch + " startdate : " + err.Error())
			continue
		}
		parsedEndDate, err := time.Parse(utils.DBDateTimeFormat, string(enddate))
		if err != nil {
			utils.LogError("Error while parsing a " + databaseTableNames.LolMatch + " enddate : " + err.Error())
			continue
		}
		lms = append(lms, LolMatch{Id: id, MatchId: matchid, Server: server, StartDate: parsedStartDate, EndDate: parsedEndDate})
	}
	return lms, nil
}

func (lm LolMatch) GetData() (*riotinterface.Match, error) {
	if !lm.dataProcessed {
		err := json.Unmarshal(lm.JsonData, &lm.data)
		if err != nil {
			utils.LogError(err.Error())
			return nil, err
		}
		lm.dataProcessed = true
	}
	return &lm.data, nil
}
