package main

import (
	"fmt"
)

import (
	"github.com/jmoiron/sqlx"
	"github.com/palestine-nights/backend/src/api"
	"github.com/palestine-nights/backend/src/db"
	"github.com/palestine-nights/backend/src/tools"
)

func initializeDB() *sqlx.DB {
	user := tools.GetEnv("DATABASE_USER", "root")
	pass := tools.GetEnv("DATABASE_PASSWORD", "")
	name := tools.GetEnv("DATABASE_NAME", "restaurant")
	host := tools.GetEnv("DATABASE_HOST", "localhost")
	port := tools.GetEnv("DATABASE_PORT", "3306")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)

	return db.Initialize(connectionString)
}

func main() {
	DB := initializeDB()
	server := api.GetServer(DB)
	server.ListenAndServe()
}
