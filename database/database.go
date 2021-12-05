package database

import (
	"database/sql"
	"les-randoms/utils"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

func OpenDatabase() {
	var err error
	Database, err = sql.Open("mysql", os.Getenv("DATABASE_CONNECTION_STRING"))
	utils.HandlePanicError(err)
	utils.LogSuccess("Database successfully opened")
}

func CloseDatabase() {
	Database.SetConnMaxIdleTime(10 * time.Second)
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
