package tools

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

const reservationSubject = "Table reservation"

type mail struct {
	senderID string
	password string
	toID     string
	subject  string
	body     string
}

type smtpServer struct {
	host string
	port string
}

// Reservation object for email.
type Reservation struct {
	Guests           int64
	Email            string
	FullName         string
	Time             time.Time
	Duration         time.Duration
	ConfirmationCode string
	CancellationCode string
}

func (s *smtpServer) serverName() string {
	return s.host + ":" + s.port
}

// FormattedTime formats time for email.
func (r *Reservation) FormattedTime() string {
	return r.Time.Format("02.01.2006, 15:04")
}

// FormattedDuration formats duration for email.
func (r *Reservation) FormattedDuration() string {
	return fmt.Sprintf("%.0f", r.Duration.Minutes())
}

func (mail *mail) buildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderID)
	message += fmt.Sprintf("To: %s\r\n", mail.toID)
	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body
	return message
}

// SendReservationEmail sends email about table reservation.
func SendReservationEmail(email string, reservation Reservation) error {
	// Construct mail object
	mail := mail{}
	mail.senderID = GetEnv("EMAIL_ADDRESS", "toursearch0@gmail.com")
	mail.password = GetEnv("EMAIL_PASSWORD", "xqauqswmqkjevksr")
	mail.toID = email
	mail.subject = reservationSubject

	// Construct message body from template
	fileName := "././templates/reservation_email.tpl"
	var textBytes bytes.Buffer
	tpl, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	err = tpl.Execute(&textBytes, &reservation)
	if err != nil {
		return err
	}
	mail.body = textBytes.String()
	messageBody := mail.buildMessage()

	smtpServer := smtpServer{host: "smtp.gmail.com", port: "465"}

	// Build auth
	auth := smtp.PlainAuth("", mail.senderID, mail.password, smtpServer.host)

	// TLS config, required for Gmail
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.serverName(), tlsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		return err
	}

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Add sender and recipient
	if err = client.Mail(mail.senderID); err != nil {
		return err
	}
	if err = client.Rcpt(mail.toID); err != nil {
		return err
	}

	// Write data
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(messageBody))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	client.Quit()

	return nil
}
