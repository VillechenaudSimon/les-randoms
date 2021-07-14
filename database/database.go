package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Test() {
	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", "217240:tD5w4$dA6$MC@tcp(mysql-villechenaud-simon.alwaysdata.net:3306)/villechenaud-simon_les-randoms")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	var testValue string
	err = db.QueryRow("SELECT name FROM BlackListItem WHERE id=1").Scan(&testValue)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("DATABASE.TEST() - END - VAR:" + testValue)
}
