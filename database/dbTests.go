package database

import (
	"errors"
	"les-randoms/utils"
)

func VerifyDatabase() {
	utils.LogInfo("Starting to verify database validity..")

	testing := utils.CreateTesting("DATABASE VALIDITY TEST")

	testDatabasePing(&testing)

	tableStates := testDatabaseTables(&testing)

	defaultAccessRightsValidity := testDefaultAccessRights(&testing)

	// Display tests conclusion
	err := testing.Conclusion()
	if err != nil {
		utils.LogError(err.Error() + " failed tests but no fatal tests failed")
		utils.LogInfo("Trying to fix errors..")
		fixTables(tableStates)
		fixDefaultAccessRightsValidity(defaultAccessRightsValidity)

		utils.LogSuccess("All fixes trials done")
		VerifyDatabase()
		return
	}
	utils.LogInfo("Program will continue..")
	return
}

func fixTables(tableStates map[string]int) {
	for key, value := range tableStates {
		if value != 0 {
			utils.LogInfo("Fixing " + key + " table..")
			if value == 2 { // table does not exists
				createTable(getSpecificTableCreationQuery(key))
			} else if value == 1 { // table exists but structure is not valid
				utils.HandlePanicError(errors.New("Fixing 'table exists but structure is not valid' error is not implemented yet. Sorry bro, GLHF !"))
			}
		}
	}
}

func fixDefaultAccessRightsValidity(b bool) {
	if !b {
		utils.LogInfo("Fixing Default AccessRights Validity..")
		DeleteDatabase("FROM AccessRight WHERE userId=1 AND path='/database'")
		AccessRight_CreateNew(1, "/database", RightTypes.Authorized)
	}
}

func testDatabasePing(testing *utils.Testing) {
	// Test to ping database
	testing.TestError(Database.Ping(), "Successful ping to database", "Failed to ping database", true)
}

func testDatabaseTables(testing *utils.Testing) map[string]int {
	// Tests on tables validity

	tableStates := make(map[string]int)
	// 0 Means exists and valid
	// 1 Means exists but structure is not valid
	// 2 Means does not exists

	for table := range tablesConfig {
		tableStates[table] = 2
	}

	/*rows, _ := SelectDatabase("name, sql FROM sqlite_schema WHERE type IN ('table','view') AND name NOT LIKE 'sqlite_%' ORDER BY 1")
	var name string
	var sql string
	for rows.Next() {
		rows.Scan(&name, &sql)
		if testing.TestStringEqual(strings.ReplaceAll(getSpecificTableCreationQuery(name), " ", ""), strings.ReplaceAll(strings.ReplaceAll(sql, "\n", " "), " ", ""), name+" table exists and is valid", name+" table exists but is not valid", false) == nil {
			tableStates[name] = 0
		} else {
			tableStates[name] = 1
		}
	}*/
	rows, err := SelectDatabase("table_name from information_schema.tables WHERE table_schema='public'")
	if err != nil {
		utils.HandlePanicError(err)
	}
	var name string
	for rows.Next() {
		rows.Scan(&name)
		if tableStates[name] != 2 {
			testing.TestBool(true, false, "if you see this, there is a problem", "Found "+name+" table but this table should not exists", true)
		} else {
			rows2, err := SelectDatabase("column_name, columns.data_type FROM information_schema.columns WHERE table_schema='public' AND table_name='" + name + "'")
			if err != nil {
				utils.HandlePanicError(err)
			}
			validColumns := 0
			var column_name string
			var data_type string
			for rows2.Next() {
				rows2.Scan(&column_name, &data_type)
				if (tablesConfig[name][column_name] == "id" && data_type == "integer") || getSQLDataType(tablesConfig[name][column_name]) == data_type {
					validColumns++
				}
			}
			if testing.TestIntEqual(len(tablesConfig[name]), validColumns, name+" table exists and has valid column data types", name+" table exists but has not valid column data types", false) == nil {
				tableStates[name] = 0
			} else {
				tableStates[name] = 1
			}
		}
	}

	for key, value := range tableStates {
		if value == 2 {
			testing.TestBool(true, false, "if you see this, there is a problem", key+" table does not exists", false)
		}
	}

	return tableStates
}

func testDefaultAccessRights(testing *utils.Testing) bool {
	right, err := AccessRight_SelectFirst("WHERE userId=1 AND path='/database'")
	b := false
	if err == nil {
		b = right.RightType == RightTypes.Authorized
	}

	testing.TestBool(true, b, "User of id 1 is authorized to access /database", "User of id 1 is forbidden from accessing /database", false)

	return b
}
