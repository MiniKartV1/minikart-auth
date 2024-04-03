package utils

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type RouteMethods struct {
	OperationId string
	APIRoute    string
	Method      string
}

func GetRoutes(router *gin.Engine, microService string) []RouteMethods {
	serverLocation := os.Getenv("SERVER_LOCATION")
	var routes []RouteMethods
	for _, route := range router.Routes() {

		if strings.HasPrefix(route.Path, "/api/"+microService) && !strings.HasPrefix(route.Path, "/api/"+microService+"/health") {
			routes = append(routes, RouteMethods{
				OperationId: extractMethodName(route.Handler),
				APIRoute:    serverLocation + route.Path,
				Method:      route.Method,
			})
		}
	}
	return routes
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
