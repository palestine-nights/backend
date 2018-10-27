package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
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

	app.initializeRouter()

	return &app
}

// Server runs application server.
func (app *App) Server(port string) {
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"})

	handler := handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(app.Router)

	server := http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Listening on port " + port)
	log.Fatal(server.ListenAndServe())
}

func (app *App) initializeRouter() {
	app.Router.HandleFunc("/tables/{id:[0-9]+}", app.getTable).Methods("GET")
	app.Router.HandleFunc("/tables/{id:[0-9]+}", app.putTable).Methods("PUT")
	app.Router.HandleFunc("/tables/{id:[0-9]+}", app.deleteTable).Methods("DELETE")

	app.Router.HandleFunc("/tables", app.postTable).Methods("POST")
	app.Router.HandleFunc("/tables", app.listTables).Methods("GET")
}

func (app *App) listTables(w http.ResponseWriter, r *http.Request) {
	table := table.getList(table{}, app.DB)
	if table == nil {
		respondWithError(w, http.StatusInternalServerError, "Error")
	} else {
		respondWithJSON(w, http.StatusOK, table)
	}
}

func (app *App) getTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table ID, must be int")
		return
	}

	table := table.getTable(table{}, app.DB, id)

	if table.ID == 0 {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Table with id %d could not be found", id))
		return
	}
	respondWithJSON(w, http.StatusOK, table)
}

func (app *App) postTable(w http.ResponseWriter, r *http.Request) {
	var table table
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&table); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if table.Places == 0 {
		respondWithError(w, http.StatusBadRequest, "Places count should be more than 0")
		return
	}
	defer r.Body.Close()
	table.ID = 0
	table.createTable(app.DB)
	respondWithJSON(w, http.StatusCreated, table)
}

func (app *App) putTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table ID, must be int")
		return
	}

	var table table

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&table); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if table.Places == 0 {
		respondWithError(w, http.StatusBadRequest, "Places count should be more than 0")
		return
	}
	defer r.Body.Close()

	table.ID = id

	// Check if id exists.
	err = table.updateTable(app.DB)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, table)
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

	// Check if id exists.
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
