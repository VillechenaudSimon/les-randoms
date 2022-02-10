package database

import (
	"database/sql"
	"errors"
	"fmt"
	"les-randoms/utils"
	"reflect"
)

type User struct {
	Id        int
	Name      string
	DiscordId string
}

func (user User) ToStringSlice() []string {
	return []string{fmt.Sprint(user.Id), user.Name, user.DiscordId}
}

func Users_ToDBStructSlice(users []User) []DBStruct {
	var r []DBStruct
	for _, user := range users {
		r = append(r, user)
	}
	return r
}

func User_GetType() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(&User{})).Type()
}

func User_SelectAll(queryPart string) ([]User, error) {
	rows, err := SelectDatabase("* FROM " + databaseTableNames.User + " " + queryPart)
	defer rows.Close()
	if err != nil {
		utils.LogError("Error while selecting on " + databaseTableNames.User + " table : " + err.Error())
		return nil, err
	}
	users := make([]User, 0)
	for rows.Next() {
		var id int
		var name string
		var discordId string
		err = rows.Scan(&id, &name, &discordId)
		if err != nil {
			utils.LogError("Error while scanning on " + databaseTableNames.User + " table : " + err.Error())
			return nil, err
		}
		users = append(users, User{Id: id, Name: name, DiscordId: discordId})
	}
	return users, nil
}

func User_SelectFirst(queryPart string) (User, error) {
	rows, err := SelectDatabase("* FROM " + databaseTableNames.User + " " + queryPart)
	if err != nil {
		utils.LogError("Error while selecting on " + databaseTableNames.User + " table : " + err.Error())
		return User{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return User{}, errors.New("No " + databaseTableNames.User + " match the request")
	}
	var id int
	var name string
	var discordId string
	err = rows.Scan(&id, &name, &discordId)
	if err != nil {
		utils.LogError("Error while scanning on " + databaseTableNames.User + " table : " + err.Error())
		return User{}, err
	}
	return User{Id: id, Name: name, DiscordId: discordId}, nil
}

func User_CreateNew(name string, discordId string) (User, sql.Result, error) {
	result, err := InsertDatabase(databaseTableNames.User + "(name, discordid) VALUES(" + utils.Esc(name) + ", " + utils.Esc(discordId) + ")")
	if err != nil {
		return User{}, result, err
	}
	newId, err := result.LastInsertId()
	if err != nil {
		return User{}, result, err
	}
	return User{Id: int(newId), Name: name, DiscordId: discordId}, result, err
}
