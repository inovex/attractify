package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"text/template"

	"attractify.io/platform/config"
	gomail "gopkg.in/gomail.v2"
)

type Mailer struct {
	config *config.Config
}

func New(config *config.Config) *Mailer {
	return &Mailer{config: config}
}

// Send sends an email to the given address
func (m Mailer) send(to, subject, body string) error {
	return m.deliver(to, subject, body, nil, "")
}

// SendWithAttachment sends an email to the given address with attachment
func (m Mailer) sendWithAttachment(to, subject, body string, attachment io.Reader, filename string) error {
	return m.deliver(to, subject, body, attachment, filename)
}

func (m Mailer) deliver(to, subject, body string, attachment io.Reader, filename string) error {
	if m.config.Environment != "production" {
		fmt.Printf("############\nTo: %s\nSubject: %s\n\n%s############", to, subject, body)
		if attachment != nil {
			fh, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				fmt.Println(err)
			}
			if _, err := io.Copy(fh, attachment); err != nil {
				fmt.Println(err)
			}
			fh.Close()
		}
		return nil
	}

	dialer := gomail.NewDialer(
		m.config.SMTP.Host,
		m.config.SMTP.Port,
		m.config.SMTP.Username,
		m.config.SMTP.Password,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", m.config.SMTP.FromAddress)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)

	if attachment != nil {
		mailer.Attach(filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, attachment)
			return err
		}))
	}

	mailer.SetBody("text/html", body)
	return dialer.DialAndSend(mailer)
}

func (m Mailer) renderTemplate(tpl string, binding interface{}) (string, error) {
	path := fmt.Sprintf("templates/mails/%s.tmpl", tpl)
	tmpl, err := template.New("layout.tmpl").ParseFiles("templates/mails/layout.tmpl", path)
	if err != nil {
		return "", err
	}
	var html bytes.Buffer
	err = tmpl.Execute(&html, binding)
	if err != nil {
		return "", err
	}
	return html.String(), nil
}
