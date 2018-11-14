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

const (
	statusCreated   = "created"
	statusApproved  = "approved"
	statusCancelled = "cancelled"
)

/* Table Reservations API */

// Create reservation.
func (server *Server) createReservation(w http.ResponseWriter, r *http.Request) {
	reservation := db.Reservation{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&reservation); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set default status "created" after creating.
	reservation.Status = statusCreated
	reservation.Duration *= time.Minute

	// Validate, that number of guests is more that 0.
	if reservation.Guests <= 0 {
		respondWithError(w, http.StatusBadRequest, "Ivalid number of guests, should more that 0")
		return
	}

	// Validate, that duration between 1h to 3h.
	if reservation.Duration < time.Hour || reservation.Duration > 3*time.Hour {
		respondWithError(w, http.StatusBadRequest, "Invalid duration time, should be between 1 and 3")
		return
	}

	isValid, err := reservation.Validate(server.DB)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate, that reservation time is not earlier than current time.
	if reservation.Time.Before(time.Now()) {
		errorMesssage := fmt.Sprintf("Invalid reservation time")
		respondWithError(w, http.StatusBadRequest, errorMesssage)
		return
	}

	// Validate, that table with TableID exists.
	if _, err = db.Table.Find(db.Table{}, server.DB, reservation.TableID); err != nil {
		errorMesssage := fmt.Sprintf("Invalid table id %d", reservation.TableID)
		respondWithError(w, http.StatusBadRequest, errorMesssage)
		return
	}

	if isValid {
		err := reservation.Insert(server.DB)
		// reservation :=

		if err != nil {
			respondWithError(w, http.StatusConflict, err.Error())
		} else {
			respondWithJSON(w, http.StatusOK, reservation)
		}
	} else {
		respondWithError(w, http.StatusConflict, "Email or phone was already used for last 24 hours")
	}
}

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
