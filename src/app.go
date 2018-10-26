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

// App is structure with router and DB instances.
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// GetApp returns applocation instances.
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
func (app *App) Server(port string) {
	log.Printf("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, app.Router))
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/tables/{id:[0-9]+}", app.getTable).Methods("GET")
	app.Router.HandleFunc("/tables/{id:[0-9]+}", app.putTable).Methods("PUT")
	app.Router.HandleFunc("/tables/{id:[0-9]+}", app.deleteTable).Methods("DELETE")

	app.Router.HandleFunc("/tables", app.postTable).Methods("POST")
	app.Router.HandleFunc("/tables", app.listTables).Methods("GET")
}

func (app *App) listTables(w http.ResponseWriter, r *http.Request) {
	t := table.getList(table{}, app.DB)
	if t == nil {
		respondWithError(w, http.StatusInternalServerError, "Error")
	} else {
		respondWithJSON(w, http.StatusOK, t)
	}
}

func (app *App) getTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table ID, must be int")
		return
	}
	t := table.getTable(table{}, app.DB, id)

	if t.ID == 0 {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Table with id %d could not be found", id))
		return
	}
	respondWithJSON(w, http.StatusOK, t)
}

func (app *App) postTable(w http.ResponseWriter, r *http.Request) {
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
	t.createTable(app.DB)
	respondWithJSON(w, http.StatusCreated, t)
}

func (app *App) putTable(w http.ResponseWriter, r *http.Request) {
	// Table id
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table ID, must be int")
		return
	}
	// Table object
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
	t.ID = id
	// Check if id exists
	err = t.updateTable(app.DB)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, t)
	}
}

func (app *App) deleteTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table ID, must be int")
		return
	}
	err = table.deleteTable(table{}, app.DB, id)
	// Check if id exists
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
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
