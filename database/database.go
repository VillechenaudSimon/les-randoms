package database

import (
	"database/sql"
	"errors"
	"les-randoms/utils"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func OpenDatabase() {
	_, err := os.Stat("sqlite-database.db")
	if err != nil { // Test if database does not exists
		utils.LogInfo("Database file missing. Creating it..")
		file, err := os.Create("sqlite-database.db") // Create SQLite file
		if err != nil {
			utils.HandlePanicError(errors.New("An error happened while creating database file : " + err.Error()))
		}
		file.Close()
		utils.LogSuccess("Database file created")
	}

	Database, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		utils.HandlePanicError(errors.New("An error happened while opening database file : " + err.Error()))
	}
	utils.LogSuccess("Database successfully opened")
}

func CloseDatabase() {
	err := Database.Close()
	utils.HandlePanicError(err)
	utils.LogSuccess("Database successfully closed")
}

func VerifyDatabase() {
	utils.LogInfo("Starting to verify database validity..")

	validTables := make(map[string]int)
	// 0 Means exists and valid
	// 1 Means exists but not valid
	// 2 Means does not exists
	validTables["User"] = 2
	validTables["Riot"] = 2

	testing := utils.CreateTesting("DATABASE VALIDITY TEST")

	// Test to ping database
	testing.TestError(Database.Ping(), "Successful ping to database", "Failed to ping database", true)

	// Tests on tables validity
	rows, _ := SelectDatabase("name, sql FROM sqlite_schema WHERE type IN ('table','view') AND name NOT LIKE 'sqlite_%' ORDER BY 1")
	var name string
	var sql string
	for rows.Next() {
		rows.Scan(&name, &sql)
		if testing.TestStringEqual(getSpecificTableCreationQuery(name), sql, name+" table exists and is valid", name+" table exists but is not valid", false) != nil {
			validTables[name] = 0
		} else {
			validTables[name] = 1
		}
	}
	for key, value := range validTables {
		if value == 2 {
			testing.TestBool(true, false, "if you see this, there is a problem", key+" table does not exists", false)
		}
	}

	// Display tests conclusion
	err := testing.Conclusion()
	if err != nil {
		utils.LogError(err.Error() + " failed tests but no fatal tests failed")
	}
	utils.LogInfo("Program will continue..")
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
