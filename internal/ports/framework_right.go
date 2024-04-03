package ports

import user_models "github.com/MiniKartV1/minikart-auth/pkg/models"

/*
	The following ports are the way to connect to the external services of the system
	which are used to hold the data and send email
*/

type DBPort interface {
	CloseDBConnection()
	AddUser(user *user_models.User) error // register
	UpdatePassword()
	UpdateLastSignedIn(email *string) (*user_models.User, error) // changepassword
	FindUserByEmail(email *string) (*user_models.User, error)
	SignOut()
	SaveCode() // creates code for the user in the database
}

type EmailServicePort interface {
	SendWelcomeEmail()
	SendResetPasswordEmail()
}
