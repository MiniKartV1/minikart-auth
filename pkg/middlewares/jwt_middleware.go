package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MiniKartV1/minikart-auth/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// jwtMiddleware checks the request for a valid JWT token
func JwtMiddleware() gin.HandlerFunc {
	// This acts as an middleware
	return func(ctx *gin.Context) {
		const Bearer_schema = "Bearer "
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, Bearer_schema) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		tokenString := authHeader[len(Bearer_schema):]
		token, err := jwt.ParseWithClaims(tokenString, &types.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secretKey := os.Getenv("SECRET_KEY")
			// Return the secret key used for signing tokens
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(*types.UserClaims)
		// Token is valid; you might want to extract claims and set them in the context
		if ok && token.Valid {
			// Attach user information to the context
			ctx.Set("user", claims)
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		ctx.Next()
	}
}
