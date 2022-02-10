package database

import (
	"errors"
	"les-randoms/utils"
)

var tablesConfig map[string]map[string]string

func init() {
	tablesConfig = map[string]map[string]string{
		databaseTableNames.User: {
			"id":        "id",
			"name":      "string",
			"discordid": "string",
		},
		databaseTableNames.List: {
			"id":      "id",
			"name":    "string",
			"headers": "text",
		},
		databaseTableNames.ListItem: {
			"id":      "id",
			"listid":  "int",
			"ownerid": "int",
			"value":   "text",
			"date":    "datetime",
		},
		databaseTableNames.AccessRight: {
			"userid":    "int",
			"path":      "string",
			"righttype": "int",
		},
		databaseTableNames.Summoner: {
			"summonerid":    "string",
			"userid":        "int",
			"accountid":     "string",
			"puuid":         "string",
			"name":          "string",
			"profileiconid": "int",
			"level":         "int",
			"revisiondate":  "int",
			"lastupdated":   "datetime",
		},
	}
}

func createTable(query string) {
	_, err := Database.Exec(query)
	if err != nil {
		utils.HandlePanicError(errors.New("SQL Query Failed while creating table : " + query + " : " + err.Error()))
	}
}

func getCreateTableQuery(tableName string) string {
	query := "CREATE TABLE " + tableName + " ( "
	i := 0
	for key, value := range tablesConfig[tableName] {
		query += key + " " + getSQLDataType(value)
		if i+1 < len(tablesConfig[tableName]) {
			query += ", "
		}
		i++
	}
	query += " )"
	return query
}

func getSQLDataType(dataType string) string {
	switch dataType {
	case "id":
		//query += "INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT"
		return "SERIAL PRIMARY KEY"
	case "int":
		//query += "INTEGER"
		return "integer"
	case "string":
		//query += "VARCHAR(32)"
		return "character varying"
	case "text":
		//query += "TEXT"
		return "text"
	case "datetime":
		//query += "DATETIME"
		return "timestamp without time zone"
	default:
		utils.LogError("Program asked for an unknown datatype : " + dataType)
		return "UNDEFINED"
	}
}

func getSpecificTableCreationQuery(tableName string) string {
	/*
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
	*/
	if tablesConfig[tableName] == nil {
		utils.HandlePanicError(errors.New("Program asked to find creation query for table '" + tableName + "' which does not exists"))
	}
	return getCreateTableQuery(tableName)
}
