package ports

import (
	"time"

	user_models "github.com/MiniKartV1/minikart-auth/pkg/models"
	user_types "github.com/MiniKartV1/minikart-auth/pkg/types"
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
	SignIn(dbUser *user_models.User, user *user_types.SigInBody) (*user_types.SignedUser, error)
	SignOut(email string, currTime *time.Time) (user_types.User, error)
	SignUp(user *user_types.SignUpBody) (*user_models.User, error)
	ResetPassword(email string) (user_types.User, error)
	ChangePassword(email, code, newPassword string) (user_types.User, error)
	GetAccessToken(token *user_models.User) (*string, error)
}
