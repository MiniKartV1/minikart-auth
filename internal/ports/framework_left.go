package ports

import "github.com/gin-gonic/gin"

/*
	This is the ports for the framework layer. here we are specifying the ports
	for our input adpater components, which can be
		1. rest api port
		2. grpc api port
*/

type RESTPort interface {
	Run()
	SignIn(ctx *gin.Context)
	SignOut(ctx *gin.Context)
	SignUp(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	Health(ctx *gin.Context)
	GetAccessToken(ctx *gin.Context)
	Profile(ctx *gin.Context)
}
