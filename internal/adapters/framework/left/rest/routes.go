package rest

import (
	"naresh/m/auth/internal/types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (rest Adapter) SignIn(ctx *gin.Context) {

	var signInBody types.SigInBody
	if err := ctx.ShouldBind(&signInBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	user, err := rest.api.SignIn(&signInBody)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "signin successful",
		"user":    user,
	})
	return
}
func (rest Adapter) SignOut(ctx *gin.Context) {
	var signObject types.SigInBody
	if err := ctx.ShouldBind(&signObject); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	res, err := rest.api.SignOut(signObject.Email)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "logout successful",
		"user":    res,
	})
	return
}
func (rest Adapter) SignUp(ctx *gin.Context) {
	var signupObject types.SignUpBody
	if err := ctx.ShouldBind(&signupObject); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}
	user, err := rest.api.SignUp(&signupObject)

	if err != nil {
		errStr := err.Error()
		str := strings.Split(errStr, ":")
		if len(str) > 0 && str[0] == "EMAIL_EXISTS" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Email Already Exists",
				"err":     err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Sign up successful.",
		"user":    user,
	})
	return
}
func (rest Adapter) ResetPassword(ctx *gin.Context) {
	var restPasswordBody types.UserEmail
	if err := ctx.ShouldBind(&restPasswordBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	user, err := rest.api.ResetPassword(&restPasswordBody)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "login successful",
		"user":    user,
	})
	return
}
func (rest Adapter) ChangePassword(ctx *gin.Context) {
	var signObject types.SigInBody
	if err := ctx.ShouldBind(&signObject); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	// user, err := rest.api.ChangePassword(signObject.Email, signObject.Codde, signObje)

	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": "Internal Server Error",
	// 	})
	// 	return
	// }

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "change password successful",
	})
	return
}