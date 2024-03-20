package ports

import (
	"naresh/m/auth/internal/models"
)

/*
	The following ports are the way to connect to the external services of the system
	which are used to hold the data and send email
*/

type DBPort interface {
	CloseDBConnection()
	AddUser(user *models.User) error // register
	UpdatePassword()
	UpdateLastSignedIn(email *string) (*models.User, error) // changepassword
	FindUserByEmail(email *string) (*models.User, error)
	SignOut()
	SaveCode() // creates code for the user in the database
}

type EmailServicePort interface {
	SendWelcomeEmail()
	SendResetPasswordEmail()
}
