package database

import (
	"database/sql"
	"les-randoms/utils"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

func OpenDatabase() {
	var err error
	Database, err = sql.Open("mysql", os.Getenv("DATABASE_CONNECTION_STRING"))
	utils.HandlePanicError(err)
}

func SelectDatabase(queryBody string) (*sql.Rows, error) {
	result, err := Database.Query("SELECT " + queryBody)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func InsertDatabase(queryBody string) (sql.Result, error) {
	result, err := Database.Exec("INSERT INTO " + queryBody)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func UpdateDatabase(queryBody string) (sql.Result, error) {
	result, err := Database.Exec("UPDATE " + queryBody)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteDatabase(queryBody string) (sql.Result, error) {
	result, err := Database.Exec("DELETE " + queryBody)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func esc(s string) string {
	return "\"" + s + "\""
}
