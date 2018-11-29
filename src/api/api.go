package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"../db"
	_ "github.com/go-sql-driver/mysql" //
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql" //
	"github.com/jmoiron/sqlx"
)

// Server is composition of router and DB instances.
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
	reservationRouter := server.Router.PathPrefix("/reservations").Subrouter()
	menuRouter := server.Router.PathPrefix("/menu").Subrouter()

	/* --- Table endpoints --- */

	tablesRouter.HandleFunc("", server.postTable).Methods("POST")
	tablesRouter.HandleFunc("", server.listTables).Methods("GET")

	tablesRouter.HandleFunc("/{id:[0-9]+}", server.getTable).Methods("GET")
	tablesRouter.HandleFunc("/{id:[0-9]+}", server.putTable).Methods("PUT")
	tablesRouter.HandleFunc("/{id:[0-9]+}", server.deleteTable).Methods("DELETE")

	/* --- Reservation endpoints --- */

	reservationRouter.HandleFunc("", server.createReservation).Methods("POST")
	reservationRouter.HandleFunc("", server.getReservations).Methods("GET")

	reservationRouter.HandleFunc("/{id:[0-9]+}", server.getReservation).Methods("GET")

	/* --- Menu endpoints --- */

	menuRouter.HandleFunc("", server.postMenuItem).Methods("POST")
	menuRouter.HandleFunc("", server.listMenu).Methods("GET")

	menuRouter.HandleFunc("/{id:[0-9]+}", server.getMenuItem).Methods("GET")
	menuRouter.HandleFunc("/{id:[0-9]+}", server.putMenuItem).Methods("PUT")
	menuRouter.HandleFunc("/{id:[0-9]+}", server.deleteMenuItem).Methods("DELETE")
}
