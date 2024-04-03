package types

import "github.com/golang-jwt/jwt/v4"

type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}
type SignedUser struct {
	User
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
type UserClaims struct {
	*User
	*jwt.RegisteredClaims
}
type UserWithPassword struct {
	User
	Password string
}

type SignUpBody struct {
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	SigInBody
}

type UserEmail struct {
	Email string `json:"email" binding:"required"`
}
type SigInBody struct {
	UserEmail
	Password string `json:"password" binding:"required"`
}
type TokenBody struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
