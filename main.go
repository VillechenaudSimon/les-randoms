package main

import (
	"les-randoms/backgroundworker"
	"les-randoms/database"
	radbwrapper "les-randoms/radb-wrapper"
	"les-randoms/webserver"
)

func main() {
	database.OpenDatabase()
	defer database.CloseDatabase()
	database.VerifyDatabase()

	go backgroundworker.Start()
	radbwrapper.SetupJobs()

	webserver.StartWebServer()
}
