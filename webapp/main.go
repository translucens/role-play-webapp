package main

import (
	"database/sql"
	"fmt"
	"log"

	"cloud.google.com/go/profiler"
	_ "github.com/googleapis/go-sql-spanner"

	"github.com/mittz/role-play-webapp/webapp/app"
	"github.com/mittz/role-play-webapp/webapp/database"
	"github.com/mittz/role-play-webapp/webapp/utils"
)

var version = "HEAD"

func main() {
	cfg := profiler.Config{
		Service:        "roleplay",
		ServiceVersion: version,
	}
	if err := profiler.Start(cfg); err != nil {
		log.Fatalln(err)
	}

	dbInfo := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		utils.GetEnvProjectID(),
		utils.GetEnvDBInstanceID(),
		utils.GetEnvDBName(),
	)

	db, err := sql.Open("spanner", dbInfo)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	dbHandler, err := database.NewDatabaseHandler("production", db)
	if err != nil {
		log.Fatal(err)
	}

	if err := dbHandler.InitDatabase(); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	const assetsDir = "./app/assets"
	const templatesDirMatch = "./app/templates/*"

	router := app.SetupRouter(dbHandler, assetsDir, templatesDirMatch)
	router.Run(":8080")
}
