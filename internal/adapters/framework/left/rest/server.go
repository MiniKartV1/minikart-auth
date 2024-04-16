package rest

import (
	"log"
	"net/http"
	"os"

	"github.com/MiniKartV1/minikart-auth/internal/ports"
	"github.com/MiniKartV1/minikart-auth/pkg/middlewares"
	user_types "github.com/MiniKartV1/minikart-auth/pkg/types"
	"github.com/MiniKartV1/minikart-auth/pkg/utils"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

type Adapter struct {
	api ports.APIPort
}

func NewAdapter(api ports.APIPort) *Adapter {
	return &Adapter{
		api: api,
	}
}

var SERVER *gin.Engine

func (rest Adapter) Run() {
	var err error
	SERVER = gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}

	// SERVER.Use(middlewares.CORSMiddleware())
	SERVER.Use(cors.New(config))
	apiRoutes := SERVER.Group("/api/auth")
	protectedRoutes := apiRoutes.Group("/protected")
	protectedRoutes.Use(middlewares.JwtMiddleware())
	registerAuthRoutes(apiRoutes, &rest)
	registerProtectedRoutes(protectedRoutes, &rest)
	port := os.Getenv("PORT")
	if len(port) != 4 {
		port = "5000"
	}
	err = SERVER.Run(":" + port)

	if err != nil {
		log.Fatalf("Cannot start the rest server %v", err)
	}
}
func (rest Adapter) Health(ctx *gin.Context) {
	if claims, exists := ctx.Get("user"); exists {
		userClaims := claims.(*user_types.UserClaims) // Type assertion
		isActive, userErr := rest.api.Health(&userClaims.Email)
		if !isActive || userErr != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "user is not active. please contact dev",
			})
		}
	}

	clientType := ctx.GetHeader("X-Client-Type")
	if clientType == "web-app" {
		authRoutes := utils.GetRoutes(SERVER, "auth")
		ctx.JSON(http.StatusOK, gin.H{"status": "UP", "operations": authRoutes})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "UP", "message": "minikart is now very light, it is exporting only bin/application"})
	return
}
func (rest Adapter) Profile(ctx *gin.Context) {
	if claims, exists := ctx.Get("user"); exists {
		userClaims := claims.(*user_types.UserClaims) // Type assertion
		// Now you can use userClaims.Email, userClaims.Roles, etc.
		ctx.JSON(http.StatusOK, gin.H{
			"email":    userClaims.Email,
			"fullname": userClaims.FirstName + " " + userClaims.LastName,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello from protected endpoint",
	})
	return
}
func registerProtectedRoutes(protectedRoutes *gin.RouterGroup, rest *Adapter) {
	protectedRoutes.GET("/profile", rest.Profile)
}
func registerAuthRoutes(apiRoutes *gin.RouterGroup, rest *Adapter) {
	apiRoutes.GET("/health", rest.Health)
	apiRoutes.POST("/signin", rest.SignIn)
	apiRoutes.POST("/signup", rest.SignUp)
	apiRoutes.POST("/signout", rest.SignOut)
	apiRoutes.POST("/reset-password", rest.ResetPassword)
	apiRoutes.POST("/change-password", rest.ChangePassword)
	apiRoutes.POST("/access-token", rest.GetAccessToken)
}

// getAuthRoutes returns a list of routes under /api/auth
