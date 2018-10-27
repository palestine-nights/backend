package main

import (
	"github.com/palestine-nights/backend/src/api"
	"github.com/palestine-nights/backend/src/tools"
)

func main() {
	databaseUser := tools.GetEnv("DATABASE_USER", "root")
	databasePassword := tools.GetEnv("DATABASE_PASSWORD", "")
	databaseName := tools.GetEnv("DATABASE_NAME", "restaurant")
	databaseHost := tools.GetEnv("DATABASE_HOST", "localhost")
	databasePort := tools.GetEnv("DATABASE_NAME", "3306")

	port := tools.GetEnv("PORT", "8080")

	server := api.GetServer(
		databaseUser,
		databasePassword,
		databaseName,
		databaseHost,
		databasePort,
	)
	server.Server(port)
}
