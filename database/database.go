package database

import (
	"database/sql"
	"les-randoms/utils"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func OpenDatabase() {
	var err error

	if _, err := os.Stat("sqlite-database.db"); err != nil { // Test if database does not exists
		utils.LogInfo("Database file missing. Creating it..")
		file, err := os.Create("sqlite-database.db") // Create SQLite file
		if err != nil {
			utils.HandlePanicError(err)
		}
		file.Close()
		utils.LogSuccess("Database file created")
	}
	Database, err = sql.Open("sqlite3", "./sqlite-database.db")

	utils.HandlePanicError(err)
	utils.LogSuccess("Database successfully opened")
}

func CloseDatabase() {
	err := Database.Close()
	utils.HandlePanicError(err)
	utils.LogSuccess("Database successfully closed")
}

func SelectDatabase(queryBody string) (*sql.Rows, error) {
	result, err := Database.Query("SELECT " + queryBody)
	if err != nil {
		utils.LogError("SQL Query Failed : SELECT " + queryBody)
		return nil, err
	}
	return result, nil
}

func InsertDatabase(queryBody string) (sql.Result, error) {
	result, err := Database.Exec("INSERT INTO " + queryBody)
	if err != nil {
		utils.LogError("SQL Query Failed : INSERT INTO " + queryBody)
		return nil, err
	}
	return result, nil
}

func UpdateDatabase(queryBody string) (sql.Result, error) {
	result, err := Database.Exec("UPDATE " + queryBody)
	if err != nil {
		utils.LogError("SQL Query Failed : UPDATE " + queryBody)
		return nil, err
	}
	return result, nil
}

func DeleteDatabase(queryBody string) (sql.Result, error) {
	result, err := Database.Exec("DELETE " + queryBody)
	if err != nil {
		utils.LogError("SQL Query Failed : DELETE " + queryBody)
		return nil, err
	}
	return result, nil
}
