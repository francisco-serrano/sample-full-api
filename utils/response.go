package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/views"
	"net/http"
)

func StatusOkResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, views.BaseResponse{
		StatusCode: http.StatusOK,
		Data:       data,
	})
}

func StatusCreatedResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, views.BaseResponse{
		StatusCode: http.StatusCreated,
		Data:       data,
	})
}

func StatusBadRequestResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusBadRequest, views.BaseResponse{
		StatusCode: http.StatusBadRequest,
		Data:       data,
	})
}

func StatusInternalErrorResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusInternalServerError, views.BaseResponse{
		StatusCode: http.StatusInternalServerError,
		Data:       data,
	})
}

func StatusServiceUnavailable(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusServiceUnavailable, views.BaseResponse{
		StatusCode: http.StatusServiceUnavailable,
		Data:       data,
	})
}
