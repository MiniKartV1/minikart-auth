package middlewares

import (
	"net/http"
	"strings"

	"github.com/MiniKartV1/minikart-auth/pkg/types"
	"github.com/MiniKartV1/minikart-auth/pkg/utils"

	"github.com/gin-gonic/gin"
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
		token, err := utils.GetUserClaimsFromToken(tokenString)

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid_Token"})
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
