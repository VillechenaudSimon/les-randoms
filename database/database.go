package database

import (
	"database/sql"
	"errors"
	"les-randoms/utils"
	"os"

	_ "github.com/lib/pq"
)

var Database *sql.DB
var databaseTableNames databaseTableNamesConst

type databaseTableNamesConst struct {
	User        string
	List        string
	ListItem    string
	AccessRight string
	Summoner    string
}

func init() {
	databaseTableNames = databaseTableNamesConst{
		User:        "appuser",
		List:        "list",
		ListItem:    "listitem",
		AccessRight: "accessright",
		Summoner:    "summoner",
	}
}

func OpenDatabase() {
	/*_, err := os.Stat("sqlite-database.db")
	if err != nil { // Test if database does not exists
		utils.LogInfo("Database file missing. Creating it..")
		file, err := os.Create("sqlite-database.db") // Create SQLite file
		if err != nil {
			utils.HandlePanicError(errors.New("An error happened while creating database file : " + err.Error()))
		}
		file.Close()
		utils.LogSuccess("Database file created")
	}*/

	var err error
	//Database, err = sql.Open("sqlite3", "./sqlite-database.db")
	utils.LogInfo("URL:   " + os.Getenv("DATABASE_URL"))
	Database, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		utils.HandlePanicError(errors.New("An error happened while opening database : " + err.Error()))
	}
	utils.LogSuccess("Database successfully opened")
}

func CloseDatabase() {
	err := Database.Close()
	utils.HandlePanicError(err)
	utils.LogSuccess("Database successfully closed")
}

func SelectDatabase(queryBody string) (*sql.Rows, error) {
	utils.LogClassic("REQUEST DATABASE : SELECT " + queryBody)
	result, err := Database.Query("SELECT " + queryBody)
	if err != nil {
		utils.LogError("SQL Query Failed : SELECT " + queryBody + " ERROR: " + err.Error())
		return nil, err
	}
	return result, nil
}

func InsertDatabase(queryBody string) (sql.Result, error) {
	utils.LogClassic("REQUEST DATABASE : INSERT INTO " + queryBody)
	result, err := Database.Exec("INSERT INTO " + queryBody)
	if err != nil {
		utils.LogError("SQL Query Failed : INSERT INTO " + queryBody + " ERROR: " + err.Error())
		return nil, err
	}
	return result, nil
}

func UpdateDatabase(queryBody string) (sql.Result, error) {
	utils.LogClassic("REQUEST DATABASE : UPDATE " + queryBody)
	result, err := Database.Exec("UPDATE " + queryBody)
	if err != nil {
		utils.LogError("SQL Query Failed : UPDATE " + queryBody + " ERROR: " + err.Error())
		return nil, err
	}
	return result, nil
}

func DeleteDatabase(queryBody string) (sql.Result, error) {
	utils.LogClassic("REQUEST DATABASE : DELETE " + queryBody)
	result, err := Database.Exec("DELETE " + queryBody)
	if err != nil {
		utils.LogError("SQL Query Failed : DELETE " + queryBody + " ERROR: " + err.Error())
		return nil, err
	}
	return result, nil
}
