package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/services"
	"github.com/sample-full-api/utils"
	"net/http"
)

type HealthController struct {
	ServiceFactory func() services.HealthService
}

// HealthCheck
// @Description Performs Health Check
// @Success 201 {object} views.BaseResponse
// @Failure 503 {object} views.BaseResponse
// @Router /forecast/health [get]
func (f *HealthController) HealthCheck(ctx *gin.Context) {
	result, errs := f.ServiceFactory().HealthCheck()
	if len(errs) != 0 {
		utils.SetResponse(ctx, http.StatusServiceUnavailable, result)
		return
	}

	utils.SetResponse(ctx, http.StatusOK, result)
}
