//Filename internal/mailer/mailer.go

package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"gopkg.in/mail.v2"
)

//go:embed "templates"
var templateFS embed.FS

// create a Mailer type
type Mailer struct {
	dailer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	dailer := mail.NewDialer(host, port, username, password)
	dailer.Timeout = 5 * time.Second
	//return our Mailer instance

	return Mailer{
		dailer: dailer,
		sender: sender,
	}
}

// send a mail
func (m Mailer) Send(recipient, templateFile string, data interface{}) error {

	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)

	if err != nil {
		return err
	}
	//Excute the template
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	//Excute the template
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	//Excute the template
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	//create new mail message

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	//call dailAndSend()
	err = m.dailer.DialAndSend(msg)

	if err != nil {
		return err
	}

	return nil
}
