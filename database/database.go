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

func SelectDatabase() {

}

func InsertDatabase() {

}

func UpdateDatabase() {

}

func DeleteDatabase() {

}
