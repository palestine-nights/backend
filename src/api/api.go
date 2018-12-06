package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql" //
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
	"github.com/jmoiron/sqlx"
	"github.com/palestine-nights/backend/src/db"
	"github.com/rs/cors"
)

// GenericError error model.
//
// swagger:model
type GenericError struct {
	// Error massage.
	Error string `json:"error"`
}

// Server is composition of router and DB instances.
// swagger:ignore
type Server struct {
	Router *mux.Router
	DB     *sqlx.DB
	DBConn *sql.Conn
}

// GetServer returns server instance.
func GetServer(user, password, database, host, port string) *Server {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)

	DB := db.Initialize(connectionString)

	router := mux.NewRouter()
	server := Server{Router: router, DB: DB}

	server.initializeRouter()

	return &server
}

func (server *Server) getHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, "./html/home.html")
}

// ListenAndServe server.
func (server *Server) ListenAndServe(port string) {
	options := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "HEAD", "OPTIONS"},
	})

	httpServer := http.Server{
		Addr:         ":" + port,
		Handler:      options.Handler(server.Router),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Listening on port " + port)
	log.Fatal(httpServer.ListenAndServe())
}

func (server *Server) initializeRouter() {
	tablesRouter := server.Router.PathPrefix("/tables").Subrouter()
	reservationRouter := server.Router.PathPrefix("/reservations").Subrouter()
	menuRouter := server.Router.PathPrefix("/menu").Subrouter()

	/* --- Homehandler --- */

	server.Router.HandleFunc("/", server.getHome).Methods("GET")

	/* --- Table endpoints --- */

	tablesRouter.HandleFunc("", server.postTable).Methods("POST")
	tablesRouter.HandleFunc("", server.listTables).Methods("GET")

	tablesRouter.HandleFunc("/{id:[0-9]+}", server.getTable).Methods("GET")
	tablesRouter.HandleFunc("/{id:[0-9]+}", server.putTable).Methods("PUT")
	tablesRouter.HandleFunc("/{id:[0-9]+}", server.deleteTable).Methods("DELETE")

	/* --- Reservation endpoints --- */

	reservationRouter.HandleFunc("", server.postReservation).Methods("POST")
	reservationRouter.HandleFunc("", server.getReservations).Methods("GET")

	reservationRouter.HandleFunc("/{id:[0-9]+}", server.getReservation).Methods("GET")
	reservationRouter.HandleFunc("/{id:[0-9]+}", server.approveReservation).Methods("POST")
	reservationRouter.HandleFunc("/{id:[0-9]+}", server.cancelReservation).Methods("POST")

	/* --- Menu endpoints --- */

	menuRouter.HandleFunc("", server.postMenuItem).Methods("POST")
	menuRouter.HandleFunc("", server.listMenu).Methods("GET")

	menuRouter.HandleFunc("/{id:[0-9]+}", server.getMenuItem).Methods("GET")
	menuRouter.HandleFunc("/{id:[0-9]+}", server.putMenuItem).Methods("PUT")
	menuRouter.HandleFunc("/{id:[0-9]+}", server.deleteMenuItem).Methods("DELETE")

	menuRouter.HandleFunc("/categories", server.getAllCategories).Methods("GET")

	// TODO: Add endpoint for categories.
	// menuRouter.HandleFunc("/categories", server.getCategories).Methods("GET")
	menuRouter.HandleFunc("/{category:[a-z|-]+}", server.listMenuByCategory).Methods("GET")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, GenericError{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
