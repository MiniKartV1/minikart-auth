package utils

import (
	"fmt"
	"os"

	user_types "github.com/MiniKartV1/minikart-auth/pkg/types"
	"github.com/golang-jwt/jwt/v4"
)

func GetUserClaimsFromToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &user_types.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secretKey := os.Getenv("SECRET_KEY")
		// Return the secret key used for signing tokens
		return []byte(secretKey), nil
	})
	return token, err
}
