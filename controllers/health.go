package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/services"
	"github.com/sample-full-api/utils"
)

type HealthController struct {
	ServiceFactory func() services.HealthService
}

// HealthCheck
// @Description Performs Health Check
// @Success 201 {object} views.BaseResponse
// @Failure 503 {object} views.BaseResponse
func (f *HealthController) HealthCheck(ctx *gin.Context) {
	result, errs := f.ServiceFactory().HealthCheck()
	if len(errs) != 0 {
		utils.StatusServiceUnavailable(ctx, result)
		return
	}

	utils.StatusOkResponse(ctx, result)
}
