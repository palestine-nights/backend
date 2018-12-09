package tools

// Mail is base representation of email.
type Mail struct {
	From    string
	To      string
	Subject string
	Body    string
}

// SMTPServer is representation of SMTP mail server.
// Created to send emails, when user cancel/confirm reservation.
type SMTPServer struct {
	Host string
	Port string
}

// URL returns full url of SMTP server.
func (s *SMTPServer) URL() string {
	return s.Host + ":" + s.Port
}
