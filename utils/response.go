package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/views"
	"net/http"
)

func StatusOkResponse(ctx *gin.Context, result interface{}) {
	ctx.JSON(http.StatusOK, views.BaseResponse{
		StatusCode: http.StatusOK,
		Data:       result,
	})
}

func StatusCreatedResponse(ctx *gin.Context, result interface{}) {
	ctx.JSON(http.StatusCreated, views.BaseResponse{
		StatusCode: http.StatusCreated,
		Data:       result,
	})
}

func StatusBadRequestResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, views.BaseResponse{
		StatusCode: http.StatusBadRequest,
		Data:       err,
	})
}
