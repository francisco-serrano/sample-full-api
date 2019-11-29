package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/services"
	"github.com/sample-full-api/utils"
	"github.com/sample-full-api/views"
	"io/ioutil"
	"strconv"
)

type ForecastController struct {
	ServiceFactory func() services.PlanetService
}

func (f *ForecastController) AddSolarSystem(ctx *gin.Context) {
	rawRequest, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		utils.StatusBadRequestResponse(ctx, err)
		return
	}

	var request views.AddSolarSystemRequest
	if err = json.Unmarshal(rawRequest, &request); err != nil {
		utils.StatusBadRequestResponse(ctx, err)
		return
	}

	result, err := f.ServiceFactory().AddSolarSystem(&request)
	if err != nil {
		utils.StatusInternalErrorResponse(ctx, err)
		return
	}

	utils.StatusCreatedResponse(ctx, result)
}

func (f *ForecastController) GetSolarSystems(ctx *gin.Context) {
	result, err := f.ServiceFactory().GetSolarSystems()
	if err != nil {
		utils.StatusInternalErrorResponse(ctx, err)
		return
	}

	utils.StatusOkResponse(ctx, result)
}

func (f *ForecastController) AddPlanet(ctx *gin.Context) {
	rawRequest, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		utils.StatusBadRequestResponse(ctx, err)
		return
	}

	var request views.AddPlanetRequest
	if err = json.Unmarshal(rawRequest, &request); err != nil {
		utils.StatusBadRequestResponse(ctx, err)
		return
	}

	result, err := f.ServiceFactory().AddPlanet(&request)
	if err != nil {
		utils.StatusInternalErrorResponse(ctx, err)
		return
	}

	utils.StatusCreatedResponse(ctx, result)
}

func (f *ForecastController) GetPlanets(ctx *gin.Context) {
	result, err := f.ServiceFactory().GetPlanets()
	if err != nil {
		utils.StatusInternalErrorResponse(ctx, err)
		return
	}

	utils.StatusOkResponse(ctx, result)
}

func (f *ForecastController) GenerateForecasts(ctx *gin.Context) {
	solarSystemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.StatusBadRequestResponse(ctx, err)
		return
	}

	daysAmount, err := strconv.Atoi(ctx.Query("days"))
	if err != nil {
		utils.StatusBadRequestResponse(ctx, err)
		return
	}

	result := f.ServiceFactory().GenerateForecasts(solarSystemId, daysAmount)

	utils.StatusOkResponse(ctx, result)
}

func (f *ForecastController) ObtainForecast(ctx *gin.Context) {
	day, err := strconv.Atoi(ctx.Query("day"))
	if err != nil {
		utils.StatusBadRequestResponse(ctx, err)
		return
	}

	result, err := f.ServiceFactory().ObtainForecast(day)
	if err != nil {
		utils.StatusInternalErrorResponse(ctx, err)
		return
	}

	utils.StatusOkResponse(ctx, result)
}
