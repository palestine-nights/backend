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

/// swagger:route POST /reservations reservations postReservation
/// Creates reservation.
/// Responses:
///   200: Reservation
func (server *Server) postReservation(w http.ResponseWriter, r *http.Request) {
	reservation := db.Reservation{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&reservation); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set default state "created" after creating.
	reservation.State = db.StateCreated

	// Validate, that number of guests is more that 0.
	if reservation.Guests <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid number of guests, should be greater that 0")
		return
	}

	// Validate, that duration between 1h to 6h.
	if reservation.Duration < time.Hour {
		respondWithError(w, http.StatusBadRequest, "Invalid duration time, should be more than "+time.Hour.String())
		return
	}
	if reservation.Duration > 6*time.Hour {
		respondWithError(w, http.StatusBadRequest, "Invalid duration time, should be less than "+(time.Hour*6).String())
		return
	}

	err := reservation.Validate(server.DB)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate, that reservation time is not earlier than current time.
	if reservation.Time.Before(time.Now()) {
		errorMessage := fmt.Sprintf("Invalid reservation time")
		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	// Validate, that table with TableID exists.
	table, err := db.Table.Find(db.Table{}, server.DB, reservation.TableID)
	if err != nil {
		errorMessage := fmt.Sprintf("Invalid table id %d", reservation.TableID)
		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}
	// Validate, that number of guests not bigger that table has.
	if reservation.Guests > table.Places {
		errorMessage := fmt.Sprintf("Invalid amount of guests, maximum amount for this table is %d", table.Places)
		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	err = reservation.Insert(server.DB)

	if err != nil {
		respondWithError(w, http.StatusConflict, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, reservation)
	}
}

/// swagger:route GET /reservations/{id} reservations getReservation
/// Returns reservation.
/// Responses:
///   200: Reservation
func (server *Server) getReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid reservation ID, must be integer")
		return
	}

	reservation, err := db.Reservation.Find(db.Reservation{}, server.DB, id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Reservation with id %d could not be found", id))
	} else {
		respondWithJSON(w, http.StatusOK, reservation)
	}
}

/// swagger:route GET /reservations/{id} reservations getReservations
/// Returns reservation.
/// Responses:
///   200: []Reservation
func (server *Server) getReservations(w http.ResponseWriter, r *http.Request) {
	getReservations := db.Reservation.GetAll

	if r.FormValue("upcoming") == "true" {
		getReservations = db.Reservation.GetUpcoming
	}

	if reservations, err := getReservations(db.Reservation{}, server.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, reservations)
	}
}

func (server *Server) updateReservationState(w http.ResponseWriter, r *http.Request, state db.State) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid reservation ID, must be integer")
		return
	}

	reservation, err := db.Reservation.Find(db.Reservation{}, server.DB, uint64(id))
	reservation.State = state

	err = reservation.Update(server.DB)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid reservation ID, must be integer")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusConflict, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, reservation.State)
	}
}

/// swagger:route POST /reservations/{id} reservations approveReservation
/// Approve reservation.
/// Responses:
///   200: State
func (server *Server) approveReservation(w http.ResponseWriter, r *http.Request) {
	server.updateReservationState(w, r, db.StateApproved)
}

/// swagger:route POST /reservations/{id} reservations cancelReservation
/// Cancel reservation.
/// Responses:
///   200: State
func (server *Server) cancelReservation(w http.ResponseWriter, r *http.Request) {
	server.updateReservationState(w, r, db.StateCancelled)
}
