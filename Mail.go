package mailer

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

// Mail struct
type Mail struct {
	to      []string
	subject string
	body    string
}

// BuildMessage of the Mail struct
func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", FROM)
	if len(mail.to) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.to, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

// Send method for Mail struct
func Send(to []string, subject, body string) {
	mail := Mail{}
	mail.to = to
	mail.subject = subject
	mail.body = body

	messageBody := mail.BuildMessage()

	log.Println(HOST)
	//build an auth
	auth := smtp.PlainAuth("", FROM, PASSWORD, HOST)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         HOST,
	}

	conn, err := tls.Dial("tcp", SERVERNAME, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, HOST)
	if err != nil {
		log.Panic(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// step 2: add all from and to
	if err = client.Mail(FROM); err != nil {
		log.Panic(err)
	}
	for _, k := range mail.to {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")
}
