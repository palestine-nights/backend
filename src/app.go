package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// App is structure with router and DB instanses.
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// GetApp returns applocation instanses.
func GetApp(user, password, database, host, port string) *App {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)

	DB, err := gorm.Open("mysql", connectionString)
	DB.AutoMigrate(&table{})

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	app := App{Router: router, DB: DB}

	app.initializeRoutes()

	return &app
}

// Server runs application server.
func (a *App) Server(port string) {
	log.Printf("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/table/{id:[0-9]+}", a.getTable).Methods("GET")
	a.Router.HandleFunc("/table", a.postTable).Methods("POST")
}

func (a *App) getTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table ID, must be int")
		return
	}
	t := table.getTable(table{}, a.DB, id)

	if t.ID == 0 {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Table with id %d could not be found", id))
		return
	}
	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) postTable(w http.ResponseWriter, r *http.Request) {
	var t table
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if t.Places == 0 {
		respondWithError(w, http.StatusBadRequest, "Places count should be more than 0")
		return
	}
	defer r.Body.Close()
	t.ID = 0
	t.createTable(a.DB)
	respondWithJSON(w, http.StatusCreated, t)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
