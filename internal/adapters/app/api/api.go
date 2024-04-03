package api

import (
	"errors"
	"time"

	"github.com/MiniKartV1/minikart-auth/internal/models"
	"github.com/MiniKartV1/minikart-auth/internal/ports"
	"github.com/MiniKartV1/minikart-auth/internal/types"
	"github.com/MiniKartV1/minikart-auth/pkg/utils"
)

/*
This is the application layer of the system This adapter has the access to ports
 1. auth or core port
 2. db port (external port)
 3. emailService port (external port)

The app layer has the dependency on domain layer which is acceptable because the outer layer depends on
the inner layer, but the dependency on the external systems is connected using the inversion of control.
*/
type Adapter struct {
	auth         ports.AuthenticationPort
	db           ports.DBPort
	emailService ports.EmailServicePort
}

func NewAdapter(auth ports.AuthenticationPort, db ports.DBPort, emailService ports.EmailServicePort) *Adapter {
	return &Adapter{
		db:           db,
		auth:         auth,
		emailService: emailService,
	}
}

func (api Adapter) SignIn(user *types.SigInBody) (*types.SignedUser, error) {
	// fetch user from the db
	dbUser, err := api.db.FindUserByEmail(&user.Email)
	res, err := api.auth.SignIn(dbUser, user)
	if err != nil {
		return &types.SignedUser{}, err
	}
	// TODO: db operation
	_, err = api.db.UpdateLastSignedIn(&user.Email)
	return res, nil
}

func (api Adapter) SignOut(email string) (types.User, error) {
	currTime := time.Now()
	res, err := api.auth.SignOut(email, &currTime)
	if err != nil {
		return types.User{}, err
	}
	return res, nil
}
func (api Adapter) SignUp(signup *types.SignUpBody) (*models.User, error) {
	/*
		check for the existing users with the same email.
		password hashing
		create account in the database
		sending confimation email
		verify

	*/
	newUser, err := api.auth.SignUp(signup)
	if err != nil {
		return &models.User{}, err
	}

	err = api.db.AddUser(newUser)
	if err != nil {
		return &models.User{}, err
	}
	return newUser, nil
}

func (api Adapter) ResetPassword(user *types.UserEmail) (types.User, error) {
	res, err := api.auth.ResetPassword(user.Email)

	if err != nil {
		return types.User{}, err
	}
	return res, nil
}

func (api Adapter) ChangePassword(email, code, newPassword string) (types.User, error) {

	res, err := api.auth.ChangePassword(email, code, newPassword)

	if err != nil {
		return types.User{}, err
	}
	return res, nil

}
func (api Adapter) Health(email *string) (bool, error) {

	dbUser, err := api.db.FindUserByEmail(email)
	if !dbUser.IsActive || err != nil {
		return false, nil
	}
	return true, nil
}
func (api Adapter) GetAccessToken(token string) (*string, error) {

	tokenClaims, err := utils.GetUserClaimsFromToken(token)
	claims, _ := tokenClaims.Claims.(*types.UserClaims)
	dbUser, err := api.db.FindUserByEmail(&claims.Email)
	if !dbUser.IsActive || err != nil {
		noToken := ""
		return &noToken, errors.New("Inactive_User")
	}
	res, err := api.auth.GetAccessToken(dbUser)
	if err != nil {
		return res, err
	}
	return res, err
}
