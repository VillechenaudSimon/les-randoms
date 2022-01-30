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
	query := "CREATE TABLE " + tableName + " (\n"
	for i := 0; i < len(columnNames); i++ {
		query += "\"" + columnNames[i] + "\" "
		switch columnTypes[i] {
		case "id":
			query += "INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT"
		case "int":
			query += "INTEGER"
		case "text":
			query += "TEXT"
		default:
			utils.HandlePanicError(errors.New("Error while creating table : Undefined column type : " + columnTypes[i]))
			return "IMPOSSIBLE TO REACH"
		}
		if i+1 < len(columnNames) {
			query += ",\n"
		}
	}
	query += "\n)"
	return query
}

func getSpecificTableCreationQuery(tableName string) string {
	switch tableName {
	case "User":
		return getCreateTableQuery("User", []string{"id", "name", "discordId"}, []string{"id", "text", "text"})
	}
	return ""
}
