package ports

import (
	"naresh/m/auth/internal/models"
	"naresh/m/auth/internal/types"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

/*
	This is the port for this domain of the system, the scope of the domain is to help user
	sign in, signout and sign up, forgot password features.

	So the requirement from this port is the adapter should implement the following functionalities
	1. Signin - takes email and password as parameters, validates and signs in the user
	2. SignOut - doesnot take any parameter and uses system time to log when the user signs out of the system
	3. Signup - takes firstname, lastname, email and password as parameters and creates new user in the system
	4. Reset Password - takes in the email and sends support email which contains code to change the password for the user who has forgotten the password
	5. Change Password - takes in new password that the user has entered and updates it in the system
*/

type AuthenticationPort interface {
	SignIn(dbUser *models.User, user *types.SigInBody) (*types.SignedUser, error)
	SignOut(email string, currTime *time.Time) (types.User, error)
	SignUp(user *types.SignUpBody) (*models.User, error)
	ResetPassword(email string) (types.User, error)
	ChangePassword(email, code, newPassword string) (types.User, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}