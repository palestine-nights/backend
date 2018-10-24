package main

func main() {
	databaseUser := getEnv("DATABASE_USER", "tours_admin")
	databasePassword := getEnv("DATABASE_PASSWORD", "ladmdetouris")
	databaseName := getEnv("DATABASE_NAME", "restaurant")
	databaseHost := getEnv("DATABASE_HOST", "localhost")
	databasePort := getEnv("DATABASE_NAME", "3306")

	port := getEnv("PORT", "8080")

	// You need to set your username and password here.
	app := GetApp(databaseUser, databasePassword, databaseName, databaseHost, databasePort)

	app.Run(port)
}
