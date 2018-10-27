package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" //
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
	"github.com/palestine-nights/backend/src/db"
)

// Server is composition of router and DB instances.
type Server struct {
	Router *mux.Router
	DB     *gorm.DB
}

// GetServer returns server instance.
func GetServer(user, password, database, host, port string) *Server {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)

	DB := db.Initialize(connectionString)

	router := mux.NewRouter()
	server := Server{Router: router, DB: DB}

	server.initializeRouter()

	return &server
}

// ListenAndServe server.
func (server *Server) ListenAndServe(port string) {
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"})

	handler := handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(server.Router)

	httpServer := http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Listening on port " + port)
	log.Fatal(httpServer.ListenAndServe())
}

func (server *Server) initializeRouter() {
	tablesRouter := server.Router.PathPrefix("/tables").Subrouter()

	tablesRouter.HandleFunc("", server.postTable).Methods("POST")
	tablesRouter.HandleFunc("", server.listTables).Methods("GET")

	tablesRouter.HandleFunc("/{id:[0-9]+}", server.getTable).Methods("GET")
	tablesRouter.HandleFunc("/{id:[0-9]+}", server.putTable).Methods("PUT")
	tablesRouter.HandleFunc("/{id:[0-9]+}", server.deleteTable).Methods("DELETE")
}

func (server *Server) listTables(w http.ResponseWriter, r *http.Request) {
	table := db.Table.GetAll(db.Table{}, server.DB)
	if table == nil {
		respondWithError(w, http.StatusInternalServerError, "Error")
	} else {
		respondWithJSON(w, http.StatusOK, table)
	}
}

func (server *Server) getTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table ID, must be int")
		return
	}

	table := db.Table.Find(db.Table{}, server.DB, id)

	if table.ID == 0 {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Table with id %d could not be found", id))
		return
	}
	respondWithJSON(w, http.StatusOK, table)
}

func (server *Server) postTable(w http.ResponseWriter, r *http.Request) {
	var table db.Table
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
	table.Save(server.DB)
	respondWithJSON(w, http.StatusCreated, table)
}

func (server *Server) putTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table ID, must be int")
		return
	}

	var table db.Table

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

	// Check if ID exists.
	err = table.Update(server.DB)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, table)
	}
}

func (server *Server) deleteTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid table ID, must be int")
		return
	}
	err = db.Table.Destroy(db.Table{}, server.DB, id)

	// Check if ID exists.
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
