package api

import (
	"fmt"
	"github.com/palestine-nights/backend/src/tools"
	"net/http"
	"strconv"
	"time"

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

	// Create reservation
	err = reservation.Insert(server.DB)

	if err != nil {
		c.JSON(http.StatusConflict, GenericError{Error: err.Error()})
		return
	}

	// Create tokens
	confirmationToken := db.Token.Generate(db.Token{}, reservation.ID, "confirm")
	cancellationToken := db.Token.Generate(db.Token{}, reservation.ID, "cancel")
	err = confirmationToken.Insert(server.DB)
	if err != nil {
		db.Reservation{}.Destroy(server.DB, reservation.ID)
		c.JSON(http.StatusBadRequest, GenericError{Error: err.Error()})
		return
	}
	err = cancellationToken.Insert(server.DB)
	if err != nil {
		db.Reservation{}.Destroy(server.DB, reservation.ID)
		c.JSON(http.StatusBadRequest, GenericError{Error: err.Error()})
		return
	}

	emailReservation := tools.Reservation{
		Guests:           reservation.Guests,
		Email:            reservation.Email,
		FullName:         reservation.FullName,
		Time:             reservation.Time,
		Duration:         reservation.Duration,
		ConfirmationCode: confirmationToken.Code,
		CancellationCode: cancellationToken.Code,
	}

	// Send email
	go tools.SendReservationEmail(reservation.Email, emailReservation)

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

/// swagger:route POST /reservations/confirm/{code}/ reservations confirmReservation
/// Confirm reservation.
/// Responses:
///   200: State
///   400: GenericError
func (server *Server) confirmReservation(c *gin.Context) {
	code := c.Param("code")
	if len(code) == 0 {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Missing reservation code"})
		return
	}

	token, err := db.Token{}.FindByCode(server.DB, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: err.Error()})
		return
	}

	if token.Type != "confirm" {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Cannot use this code to confirm reservation"})
		return
	}

	if token.State == "used" {
		c.JSON(http.StatusBadRequest, GenericError{Error: "This code is already used"})
		return
	}

	reservation, err := db.Reservation{}.Find(server.DB, token.ReservationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: err.Error()})
		return
	}

	if reservation.State != "created" {
		message := "Cannot confirm reservation because it has state: " + string(reservation.State)
		c.JSON(http.StatusBadRequest, GenericError{Error: message})
		return
	}

	reservation.State = "confirmed"
	err = reservation.Update(server.DB)

	token.State = "used"
	token.UpdateState(server.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}

/// swagger:route POST /reservations/approve/{id} reservations cancelReservation
/// Cancel reservation.
/// Responses:
///   200: State
///   400: GenericError
func (server *Server) cancelReservation(c *gin.Context) {
	// Cancel by id
	_, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err == nil {
		server.updateReservationState(c, db.StateCancelled)
		return
	}

	// Cancel by code
	code := c.Param("id")
	if len(code) == 0 {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Missing reservation code"})
		return
	}

	token, err := db.Token{}.FindByCode(server.DB, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: err.Error()})
		return
	}

	if token.Type != "cancel" {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Cannot use this code to cancel reservation"})
		return
	}

	if token.State == "used" {
		c.JSON(http.StatusBadRequest, GenericError{Error: "This code is already used"})
		return
	}

	reservation, err := db.Reservation{}.Find(server.DB, token.ReservationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: err.Error()})
		return
	}

	reservation.State = "cancelled"
	err = reservation.Update(server.DB)

	token.State = "used"
	token.UpdateState(server.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}
