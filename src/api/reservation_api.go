package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/palestine-nights/backend/src/db"
)

/* Table Reservations API */

/// swagger:route POST /reservations reservations postReservation
/// Creates reservation.
/// Responses:
///   200: Reservation
///   400: GenericError
func (server *Server) postReservation(c *gin.Context) {
	reservation := db.Reservation{}

	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid request payload"})
		return
	}

	// Set default state "created" after creating.
	reservation.State = db.StateCreated

	// Validate, that number of guests is more that 0.
	if reservation.Guests <= 0 {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid number of guests, should be greater that 0"})
		return
	}

	// Validate, that duration between 1h to 6h.
	if reservation.Duration < time.Hour {
		errorMsg := fmt.Sprintf("Invalid duration time, should be more than %s", time.Hour.String())
		c.JSON(http.StatusBadRequest, GenericError{Error: errorMsg})
		return
	}
	if reservation.Duration > 6*time.Hour {
		errorMsg := fmt.Sprintf("Invalid duration time, should be less than %s", time.Hour.String())
		c.JSON(http.StatusBadRequest, GenericError{Error: errorMsg})
		return
	}

	err := reservation.Validate(server.DB)

	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: err.Error()})
		return
	}

	// Validate, that reservation time is not earlier than current time.
	if reservation.Time.Before(time.Now()) {
		errorMsg := fmt.Sprintf("Invalid reservation time")
		c.JSON(http.StatusBadRequest, GenericError{Error: errorMsg})
		return
	}

	// Validate, that table with TableID exists.
	table, err := db.Table.Find(db.Table{}, server.DB, reservation.TableID)
	if err != nil {
		errorMsg := fmt.Sprintf("Invalid table id %d", reservation.TableID)
		c.JSON(http.StatusBadRequest, GenericError{Error: errorMsg})
		return
	}
	// Validate, that number of guests not bigger that table has.
	if reservation.Guests > table.Places {
		errorMsg := fmt.Sprintf("Invalid amount of guests, maximum amount for this table is %d", table.Places)
		c.JSON(http.StatusBadRequest, GenericError{Error: errorMsg})
		return
	}

	err = reservation.Insert(server.DB)

	if err == nil {
		c.JSON(http.StatusOK, reservation)
	} else {
		c.JSON(http.StatusConflict, GenericError{Error: err.Error()})
	}
}

/// swagger:route GET /reservations/{id} reservations getReservation
/// Returns reservation.
/// Responses:
///   200: Reservation
///   404: GenericError
func (server *Server) getReservation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid reservation ID, must be integer"})
		return
	}

	reservation, err := db.Reservation.Find(db.Reservation{}, server.DB, id)

	if err == nil {
		c.JSON(http.StatusOK, reservation)
	} else {
		errorMsg := fmt.Sprintf("Reservation with id %d could not be found", id)
		c.JSON(http.StatusNotFound, GenericError{Error: errorMsg})
	}
}

/// swagger:route GET /reservations/{id} reservations getReservations
/// Returns reservation.
/// Responses:
///   200: []Reservation
///   500: []Reservation
func (server *Server) getReservations(c *gin.Context) {
	if reservations, err := db.Reservation.GetAll(db.Reservation{}, server.DB); err == nil {
		c.JSON(http.StatusOK, reservations)
	} else {
		c.JSON(http.StatusInternalServerError, GenericError{Error: err.Error()})
	}
}

func (server *Server) updateReservationState(c *gin.Context, state db.State) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Invalid reservation ID, must be integer"})
		return
	}

	reservation, err := db.Reservation.Find(db.Reservation{}, server.DB, uint64(id))

	if err != nil {
		errorMsg := fmt.Sprintf("Reservation with ID %d does not exist", id)
		c.JSON(http.StatusBadRequest, GenericError{Error: errorMsg})
		return
	}

	reservation.State = state
	err = reservation.Update(server.DB)

	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: err.Error()})
		return
	}

	if err == nil {
		c.JSON(http.StatusOK, reservation)
	} else {
		c.JSON(http.StatusBadRequest, GenericError{Error: err.Error()})
	}
}

/// swagger:route POST /reservations/cancel/{id} reservations approveReservation
/// Approve reservation.
/// Responses:
///   200: State
///   400: GenericError
func (server *Server) approveReservation(c *gin.Context) {
	server.updateReservationState(c, db.StateApproved)
}

/// swagger:route POST /reservations/approve/{id} reservations cancelReservation
/// Cancel reservation.
/// Responses:
///   200: State
///   400: GenericError
func (server *Server) cancelReservation(c *gin.Context) {
	server.updateReservationState(c, db.StateCancelled)
}
