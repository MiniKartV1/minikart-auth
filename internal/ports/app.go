package ports

import (
	"github.com/MiniKartV1/minikart-auth/internal/models"
	"github.com/MiniKartV1/minikart-auth/internal/types"
)

/*
	This layer sits between framework and domain layer. This layer orchestrates the request
	This layer is connected to the core of the current system and the also connected
	to the external entities of the system such as db and email service etc.
*/

type APIPort interface {
	SignIn(*types.SigInBody) (*types.SignedUser, error)
	SignOut(email string) (types.User, error)
	SignUp(*types.SignUpBody) (*models.User, error)
	ResetPassword(*types.UserEmail) (types.User, error)
	ChangePassword(email, code, newPassword string) (types.User, error)
	Health(email *string) (bool, error)
	GetAccessToken(token string) (*string, error)
}
