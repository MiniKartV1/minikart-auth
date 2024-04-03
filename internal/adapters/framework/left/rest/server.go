package rest

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MiniKartV1/minikart-auth/internal/ports"
	"github.com/MiniKartV1/minikart-auth/internal/types"
	"github.com/MiniKartV1/minikart-auth/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type RouteMethods struct {
	OperationId string
	APIRoute    string
	Method      string
}
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
	SERVER.Use(middlewares.CORSMiddleware())
	apiRoutes := SERVER.Group("/api/auth")
	protectedRoutes := apiRoutes.Group("/protected")
	protectedRoutes.Use(middlewares.JwtMiddleware())
	registerAuthRoutes(apiRoutes, &rest)
	registerProtectedRoutes(protectedRoutes)
	err = SERVER.Run(":3000")

	if err != nil {
		log.Fatalf("Cannot start the rest server %v", err)
	}
}
func (rest Adapter) Health(ctx *gin.Context) {
	if claims, exists := ctx.Get("user"); exists {
		userClaims := claims.(*types.UserClaims) // Type assertion
		isActive, userErr := rest.api.Health(&userClaims.Email)
		if !isActive || userErr != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "user is not active. please contact dev",
			})
		}
	}

	clientType := ctx.GetHeader("X-Client-Type")
	if clientType == "web-app" {
		authRoutes := getAuthRoutes(SERVER)
		ctx.JSON(http.StatusOK, gin.H{"status": "UP", "operations": authRoutes})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
	return
}
func registerProtectedRoutes(protectedRoutes *gin.RouterGroup) {
	protectedRoutes.GET("/profile", func(ctx *gin.Context) {
		if claims, exists := ctx.Get("user"); exists {
			userClaims := claims.(*types.UserClaims) // Type assertion
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
	})
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
func getAuthRoutes(router *gin.Engine) []RouteMethods {
	serverLocation := os.Getenv("SERVER_LOCATION")
	var authRoutes []RouteMethods
	for _, route := range router.Routes() {

		if strings.HasPrefix(route.Path, "/api/auth") && !strings.HasPrefix(route.Path, "/api/auth/health") {
			authRoutes = append(authRoutes, RouteMethods{
				OperationId: extractMethodName(route.Handler),
				APIRoute:    serverLocation + route.Path,
				Method:      route.Method,
			})
		}
	}
	return authRoutes
}

func extractMethodName(s string) string {
	// example string; we are extracting the method name from the below string.
	// "github.com/MiniKartV1/minikart-auth/internal/adapters/framework/left/rest.Adapter.SignIn-fm"
	// Find the last index of "/"
	lastDotIndex := strings.LastIndex(s, ".")
	// Extract the substring up to this position
	if lastDotIndex != -1 {
		return s[lastDotIndex+1 : len(s)-3]
	}
	return ""
}
