package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/smtp"
	"strconv"
	"time"
)
import (
	"github.com/alecthomas/template"
	"github.com/gin-gonic/gin"
	"github.com/palestine-nights/backend/src/db"
	"github.com/palestine-nights/backend/src/tools"
)

// ReservationEmail is Reservation struct with confirm/cancel codes.
type ReservationEmail struct {
	db.Reservation
	ConfirmationCode string
	CancellationCode string
}

// SendReservationEmail sends email about table reservation.
func SendReservationEmail(reservation ReservationEmail) error {
	// Construct mail object.

	fileName := "./templates/reservation_email.tpl"

	textBytes := bytes.Buffer{}

	tpl := template.Must(template.ParseFiles(fileName))

	err := tpl.Execute(&textBytes, reservation)

	if err != nil {
		return err
	}

	mail := tools.Mail{
		From:    tools.GetEnv("SENDER_EMAIL", "noreply@palestinenights.com"),
		To:      reservation.Email,
		Subject: "Table Reservation",
		Body:    textBytes.String(),
	}

	// Initialize mail server instance.
	mailServer := tools.SMTPServer{
		Host: tools.GetEnv("SMTP_HOST", "0.0.0.0"),
		Port: tools.GetEnv("SMTP_PORT", "1025"),
	}

	conn, err := smtp.Dial(mailServer.URL())

	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	conn.Mail(mail.From)
	conn.Rcpt(mail.To)

	wc, err := conn.Data()

	if err != nil {
		fmt.Println(err)
	}

	defer wc.Close()

	buf := bytes.NewBufferString(mail.Body)

	if _, err = buf.WriteTo(wc); err != nil {
		fmt.Println(err)
	}

	return nil
}

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

	// Create confirm token.
	err = db.GenerateToken(reservation.ID, db.TypeConfirm, server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericError{Error: err.Error()})
		return
	}

	// Create cancel token.
	err = db.GenerateToken(reservation.ID, db.TypeCancel, server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericError{Error: err.Error()})
		return
	}

	cancellationToken := db.MustLastToken(server.DB, db.TypeCancel)
	confirmationToken := db.MustLastToken(server.DB, db.TypeConfirm)

	// Send email
	err = SendReservationEmail(ReservationEmail{
		reservation,
		confirmationToken.Code,
		cancellationToken.Code,
	})

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

/// swagger:route POST /reservations/approve/{id} reservations approveReservation
/// Approve reservation.
/// Responses:
///   200: State
///   400: GenericError
func (server *Server) approveReservation(c *gin.Context) {
	server.updateReservationState(c, db.StateApproved)
}

/// swagger:route POST /reservations/cancel/{id} reservations cancelReservation
/// Cancel reservation.
/// Responses:
///   200: State
///   400: GenericError
func (server *Server) cancelReservation(c *gin.Context) {
	server.updateReservationState(c, db.StateCancelled)
}

/// swagger:route GET /confirm/{code} reservations confirmReservation
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

	if token.Type != db.TypeConfirm {
		c.JSON(http.StatusBadRequest, GenericError{Error: "Cannot use this code to confirm reservation"})
		return
	}

	if token.Used {
		c.JSON(http.StatusBadRequest, GenericError{Error: "This code is already used"})
		return
	}

	reservation, err := db.Reservation{}.Find(server.DB, token.ReservationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericError{Error: err.Error()})
		return
	}

	if reservation.State != db.StateCreated {
		message := "Cannot confirm reservation because it has state: " + string(reservation.State)
		c.JSON(http.StatusBadRequest, GenericError{Error: message})
		return
	}

	reservation.State = db.StateConfirmed
	err = reservation.Update(server.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericError{Error: err.Error()})
		return
	}

	token.Used = true
	err = token.UpdateState(server.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, reservation.State)
	}
}
