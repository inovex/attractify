package mailer

import "io"

type userRegisterBinding struct {
	RegistrationLink string
	InviterEmail     string
}

type userPasswordResetBinding struct {
	Email     string
	ResetLink string
}

func (m Mailer) SendRegistration(email, inviterEmail, link string) error {
	binding := &userRegisterBinding{RegistrationLink: link, InviterEmail: inviterEmail}
	body, err := m.renderTemplate("registration", binding)
	if err != nil {
		return err
	}
	return m.send(
		email,
		"Welcome to Attractify",
		body,
	)
}

func (m Mailer) SendPasswordReset(email, link string) error {
	binding := &userPasswordResetBinding{Email: email, ResetLink: link}
	body, err := m.renderTemplate("reset_password", binding)
	if err != nil {
		return err
	}
	return m.send(
		email,
		"Your password reset",
		body,
	)
}

func (m Mailer) SendDataExport(email string, export io.Reader) error {
	body, err := m.renderTemplate("data_export", nil)
	if err != nil {
		return err
	}
	return m.sendWithAttachment(
		email,
		"Your data export",
		body,
		export,
		"export.json",
	)
}
