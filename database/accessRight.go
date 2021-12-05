package database

import (
	"errors"
	"fmt"
	"les-randoms/utils"
	"reflect"
)

type AccessRight struct {
	UserId    int
	Path      string
	RightType bool
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
		var rightType bool
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
	defer rows.Close()
	if err != nil {
		utils.LogError("Error while selecting on AccessRight table : " + err.Error())
		return AccessRight{}, err
	}
	if !rows.Next() {
		return AccessRight{}, errors.New("No AccessRight match the request")
	}
	var userId int
	var path string
	var rightType bool
	err = rows.Scan(&userId, &path, &rightType)
	if err != nil {
		utils.LogError("Error while scanning on AccessRight table : " + err.Error())
		return AccessRight{}, err
	}
	return AccessRight{UserId: userId, Path: path, RightType: rightType}, nil
}
