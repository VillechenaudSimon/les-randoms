package database

import (
	"database/sql"
	"fmt"
	"les-randoms/utils"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

func OpenDatabase() {
	Database, err := sql.Open("mysql", os.Getenv("DATABASE_CONNECTION_STRING"))
	utils.HandlePanicError(err)

	var testValue string
	err = Database.QueryRow("SELECT name FROM BlackListItem WHERE id=1").Scan(&testValue)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("DATABASE.TEST() - END - VAR:" + testValue)
}

func SelectDatabase(queryBody string) (*sql.Rows, error) {
	result, err := Database.Query("SELECT " + queryBody)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func InsertDatabase(queryBody string) (sql.Result, error) {
	result, err := Database.Exec("INSERT" + queryBody)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func UpdateDatabase(queryBody string) (sql.Result, error) {
	result, err := Database.Exec("UPDATE" + queryBody)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteDatabase(queryBody string) (sql.Result, error) {
	result, err := Database.Exec("DELETE" + queryBody)
	if err != nil {
		return nil, err
	}
	return result, nil
}
