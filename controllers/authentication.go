package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/services"
	"github.com/sample-full-api/utils"
	"github.com/sample-full-api/views"
	"net/http"
)

type AuthenticationController struct {
	AuthServiceFactory func() services.AuthenticationService
}

// Login
// @Description Obtains token for application use
// @Success 200 {object} views.BaseResponse
// @Failure 400 {object} views.BaseResponse
// @Failure 401 {object} views.BaseResponse
// @Router /authentication/login [post]
func (a *AuthenticationController) Login(ctx *gin.Context) {
	var request views.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.SetResponse(ctx, http.StatusBadRequest, err)
		return
	}

	token, err := a.AuthServiceFactory().GenerateToken(request)
	if err != nil {
		utils.SetResponse(ctx, http.StatusUnauthorized, err)
		return
	}

	utils.SetResponse(ctx, http.StatusOK, token)
}
