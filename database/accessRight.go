package database

import (
	"database/sql"
	"errors"
	"fmt"
	"les-randoms/utils"
	"reflect"
)

var RightTypes RightTypesConst

func init() {
	RightTypes = RightTypesConst{
		Hidden:     -1,
		Forbidden:  0,
		Authorized: 1,
	}
}

type RightTypesConst struct {
	Hidden     int
	Forbidden  int
	Authorized int
}

type AccessRight struct {
	UserId    int
	Path      string
	RightType int
}

func (accessRight AccessRight) ToStringSlice() []string {
	return []string{fmt.Sprint(accessRight.UserId), accessRight.Path, fmt.Sprint(accessRight.RightType)}
}

func AccessRights_ToDBStructSlice(accessRights []AccessRight) []DBStruct {
	var r []DBStruct
	for _, accessRight := range accessRights {
		r = append(r, accessRight)
	}
	return r
}

func AccessRight_GetType() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(&AccessRight{})).Type()
}

func AccessRight_SelectAll(queryPart string) ([]AccessRight, error) {
	rows, err := SelectDatabase("userId, path, rightType FROM AccessRight " + queryPart)
	defer rows.Close()
	if err != nil {
		utils.LogError("Error while selecting on AccessRight table : " + err.Error())
		return nil, err
	}
	accessRights := make([]AccessRight, 0)
	for rows.Next() {
		var userId int
		var path string
		var rightType int
		err = rows.Scan(&userId, &path, &rightType)
		if err != nil {
			utils.LogError("Error while scanning on AccessRight table : " + err.Error())
			return nil, err
		}
		accessRights = append(accessRights, AccessRight{UserId: userId, Path: path, RightType: rightType})
	}
	return accessRights, nil
}

func AccessRight_SelectFirst(queryPart string) (AccessRight, error) {
	rows, err := SelectDatabase("userId, path, rightType FROM AccessRight " + queryPart)
	if err != nil {
		utils.LogError("Error while selecting on AccessRight table : " + err.Error())
		return AccessRight{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return AccessRight{}, errors.New("No AccessRight match the request")
	}
	var userId int
	var path string
	var rightType int
	err = rows.Scan(&userId, &path, &rightType)
	if err != nil {
		utils.LogError("Error while scanning on AccessRight table : " + err.Error())
		return AccessRight{}, err
	}
	return AccessRight{UserId: userId, Path: path, RightType: rightType}, nil
}

func AccessRight_CreateNew(userId int, path string, rightType int) (AccessRight, sql.Result, error) {
	result, err := InsertDatabase("AccessRight(userId, path, rightType) VALUES(" + fmt.Sprint(userId) + ", " + utils.Esc(path) + ", " + fmt.Sprint(rightType) + ")")
	if err != nil {
		return AccessRight{}, result, err
	}
	return AccessRight{UserId: userId, Path: path, RightType: rightType}, result, err
}
