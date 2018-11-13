package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/palestine-nights/backend/src/db"
)

/* Table Reservations API */

// Create reservation.
func (server *Server) createReservation(w http.ResponseWriter, r *http.Request) {
	var reservation db.Reservation

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reservation); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validation should be here. Email uniqueness, time uniqueness... etc.
	// var reservations []db.Reservation
	// sql := `SELECT * FROM reservations WHERE (created_at >= NOW() - INTERVAL 1 DAY) AND (email != :email OR phone != :phone)`

	// result, err := db.NamedExec(sql, reservation)

	isValid, err := reservation.Validate(server.DB)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if isValid {
		// Insert
	} else {
		respondWithError(w, http.StatusConflict, "email or phone was used for last 24 hours.")
	}

	// TODO: Add validation on email (shouldn't be used at last 24h).
	// TODO: Add validation on phone (shouldn't be used at last 24h).
	// TODO: Check that this table at this time is free.
	// TODO: Check that duration should be between 1h to 4h.

	// Check if table is available.

	// defer r.Body.Close()

	respondWithJSON(w, http.StatusCreated, reservation)
}

func (server *Server) getReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid reservation ID, must be int")
		return
	}

	reservation, err := db.Reservation.Find(db.Reservation{}, server.DB, id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Reservation with id %d could not be found", id))
		return
	}

	respondWithJSON(w, http.StatusOK, reservation)
}

func (server *Server) getAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, _ := db.Reservation.GetAll(db.Reservation{}, server.DB)

	respondWithJSON(w, http.StatusOK, reservations)
}

func (server *Server) getReservations(w http.ResponseWriter, r *http.Request) {
	reservations, _ := db.Reservation.Where(db.Reservation{}, server.DB, "time >= ?", time.Now)

	respondWithJSON(w, http.StatusOK, reservations)
}
