package database

import (
	"errors"
	"les-randoms/utils"
)

func createTable(query string) {
	_, err := Database.Exec(query)
	if err != nil {
		utils.HandlePanicError(errors.New("SQL Query Failed while creating table : " + query))
	}
}

func getCreateTableQuery(tableName string, columnNames []string, columnTypes []string) string {
	query := "CREATE TABLE " + tableName + " ( "
	for i := 0; i < len(columnNames); i++ {
		query += "\"" + columnNames[i] + "\" "
		switch columnTypes[i] {
		case "id":
			query += "INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT"
		case "int":
			query += "INTEGER"
		case "string":
			query += "VARCHAR(32)"
		case "text":
			query += "TEXT"
		case "datetime":
			query += "DATETIME"
		default:
			utils.HandlePanicError(errors.New("Error while creating table : Undefined column type : " + columnTypes[i]))
			return "IMPOSSIBLE TO REACH"
		}
		if i+1 < len(columnNames) {
			query += ", "
		}
	}
	query += " )"
	return query
}

func getSpecificTableCreationQuery(tableName string) string {
	switch tableName {
	case "User":
		return getCreateTableQuery("User", []string{"id", "name", "discordId"}, []string{"id", "string", "string"})
	case "List":
		return getCreateTableQuery("List", []string{"id", "name", "headers"}, []string{"id", "string", "text"})
	case "ListItem":
		return getCreateTableQuery("ListItem", []string{"id", "listId", "ownerId", "value", "date"}, []string{"id", "int", "int", "text", "datetime"})
	case "AccessRight":
		return getCreateTableQuery("AccessRight", []string{"userId", "path", "rightType"}, []string{"int", "string", "int"})
	case "Summoner":
		return getCreateTableQuery("Summoner", []string{"summonerId", "userId", "accountId", "puuid", "name", "profileIconId", "level", "revisionDate", "lastUpdated"}, []string{"string", "int", "string", "string", "string", "int", "int", "int", "datetime"})
	}
	return ""
}
