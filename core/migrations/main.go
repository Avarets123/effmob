package main

import (
	"effect-mobile/pkg/logger"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/joho/godotenv"
)

func main() {
	logger := logger.GetLogger()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Error(".env file not found!")
	}

	fmt.Println(os.Environ())

	// var migrator *migrate.Migrate

	// // if migrationHost == "DOCKER" {
	// // 	pwd, _ := os.Getwd()
	// // 	migrator, err = migrate.New("file:"+pwd+"/"+"migrations/files", os.Getenv("POSTGRES_URL"))
	// // } else {
	// // }

	migrator, err := migrate.New("file://migrations/files", os.Getenv("POSTGRES_URL"))
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	err = migrator.Up()

	if err != nil {

		isNoChanges := strings.Contains(err.Error(), "no change")

		if isNoChanges {
			logger.Info("There are no changes in migrations")
			return
		}

		logger.Error(err.Error())
		panic(err)

	}

}
