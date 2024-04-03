package ports

import (
	user_models "github.com/MiniKartV1/minikart-auth/pkg/models"
	user_types "github.com/MiniKartV1/minikart-auth/pkg/types"
)

/*
	This layer sits between framework and domain layer. This layer orchestrates the request
	This layer is connected to the core of the current system and the also connected
	to the external entities of the system such as db and email service etc.
*/

type APIPort interface {
	SignIn(*user_types.SigInBody) (*user_types.SignedUser, error)
	SignOut(email string) (user_types.User, error)
	SignUp(*user_types.SignUpBody) (*user_models.User, error)
	ResetPassword(*user_types.UserEmail) (user_types.User, error)
	ChangePassword(email, code, newPassword string) (user_types.User, error)
	Health(email *string) (bool, error)
	GetAccessToken(token string) (*string, error)
}
