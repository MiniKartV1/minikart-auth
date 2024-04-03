package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/MiniKartV1/minikart-auth/internal/models"
	"github.com/MiniKartV1/minikart-auth/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Adapter struct{}

func NewAdapter() *Adapter {

	return &Adapter{}
}

func (auth Adapter) SignIn(dbUser *models.User, user *types.SigInBody) (*types.SignedUser, error) {
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		return nil, errors.New("CREDENTIALS_FAILED: email or password is wrong.") // Fixed typo
	}
	accessTokenSecret := os.Getenv("SECRET_KEY")
	refreshTokenSecret := os.Getenv("REFRESH_SECRET_KEY")
	accessTokenExpiryStr := os.Getenv("ACCESS_TOKEN_EXPIRY")
	accessTokenExpiry, err := strconv.Atoi(accessTokenExpiryStr)

	accessToken, err := generateTokenForUser(
		dbUser,
		accessTokenSecret,
		time.Now().Add(time.Hour*time.Duration(accessTokenExpiry)),
	)
	refreshToken, err := generateTokenForUser(
		dbUser,
		refreshTokenSecret,
		time.Now().Add(24*time.Hour*30),
	)
	if err != nil {
		return nil, err
	}
	fmt.Println("Login Successful")
	return &types.SignedUser{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: types.User{
			FirstName: dbUser.FirstName,
			LastName:  dbUser.LastName,
			Email:     dbUser.Email,
		},
	}, nil

}

func (auth Adapter) SignUp(user *types.SignUpBody) (*models.User, error) {

	// hash the password here
	hashedPassword, err := generatePasswordHash(user.Password)
	if err != nil {
		fmt.Println("Error in generating password hash", err)
		return nil, err
	}
	newUser := models.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	}
	return &newUser, nil
}
func (auth Adapter) SignOut(email string, currTime *time.Time) (types.User, error) {
	fmt.Println(currTime)
	return types.User{}, nil
}
func (auth Adapter) ResetPassword(email string) (types.User, error) {
	return types.User{}, nil
}
func (auth Adapter) ChangePassword(email, code, newPassword string) (types.User, error) {
	return types.User{}, nil
}
func (auth Adapter) GetAccessToken(user *models.User) (*string, error) {
	accessToken := ""
	accessTokenSecret := os.Getenv("SECRET_KEY")
	accessTokenExpiryStr := os.Getenv("ACCESS_TOKEN_EXPIRY")
	accessTokenExpiry, _ := strconv.Atoi(accessTokenExpiryStr)
	accessToken, err := generateTokenForUser(
		user,
		accessTokenSecret,
		time.Now().Add(time.Hour*time.Duration(accessTokenExpiry)),
	)
	if err != nil {
		return &accessToken, err
	}
	return &accessToken, nil
}
func generatePasswordHash(password string) (string, error) {
	// Generate a hashed version of the password
	// The second argument is the cost of hashing, which determines how much time is needed to calculate the hash.
	// The higher the cost, the more secure and slow to generate the hash. 10 is a reasonable default.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func generateTokenForUser(user *models.User, secretKey string, expiresAt time.Time) (string, error) {

	claims := &types.UserClaims{
		User: &types.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt), // token expires after 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}
	return signedToken, nil
}
