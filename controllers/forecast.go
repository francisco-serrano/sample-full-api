package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/services"
	"github.com/sample-full-api/utils"
	"github.com/sample-full-api/views"
	"net/http"
	"strconv"
)

type ForecastController struct {
	ServiceFactory func() services.ForecastService
}

// AddSolarSystem
// @Description Adds a solar system into the database
// @Param message body views.AddSolarSystemRequest false "Add Solar System Request"
// @Success 201 {object} views.BaseResponse
// @Failure 400 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /forecast/solar_systems [post]
func (f *ForecastController) AddSolarSystem(ctx *gin.Context) {
	var request views.AddSolarSystemRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.SetResponse(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := f.ServiceFactory().AddSolarSystem(&request)
	if err != nil {
		utils.SetResponse(ctx, utils.GetStatusErrorCode(err), err)
		return
	}

	utils.SetResponse(ctx, http.StatusCreated, result)
}

// GetSolarSystems
// @Description Gets all added solar systems
// @Success 201 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /forecast/solar_systems [get]
func (f *ForecastController) GetSolarSystems(ctx *gin.Context) {
	result, err := f.ServiceFactory().GetSolarSystems()
	if err != nil {
		utils.SetResponse(ctx, utils.GetStatusErrorCode(err), err)
		return
	}

	utils.SetResponse(ctx, http.StatusOK, result)
}

// AddPlanet
// @Description Adds a planet into the database
// @Param message body views.AddPlanetRequest false "Add Planet Request"
// @Success 201 {object} views.BaseResponse
// @Failure 400 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /forecast/planets [post]
func (f *ForecastController) AddPlanet(ctx *gin.Context) {
	var request views.AddPlanetRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.SetResponse(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := f.ServiceFactory().AddPlanet(&request)
	if err != nil {
		utils.SetResponse(ctx, utils.GetStatusErrorCode(err), err)
		return
	}

	utils.SetResponse(ctx, http.StatusCreated, result)
}

// GetPlanets
// @Description Gets all added planets
// @Success 201 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /forecast/planets [get]
func (f *ForecastController) GetPlanets(ctx *gin.Context) {
	result, err := f.ServiceFactory().GetPlanets()
	if err != nil {
		utils.SetResponse(ctx, utils.GetStatusErrorCode(err), err)
		return
	}

	utils.SetResponse(ctx, http.StatusOK, result)
}

// GenerateForecasts
// @Description Generates forecasts given added system and planets
// @Success 201 {object} views.BaseResponse
// @Failure 400 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /forecast/solar_systems/{id}/generate_forecasts [post]
func (f *ForecastController) GenerateForecasts(ctx *gin.Context) {
	solarSystemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.SetResponse(ctx, http.StatusBadRequest, err)
		return
	}

	daysAmount, err := strconv.Atoi(ctx.Query("days"))
	if err != nil {
		utils.SetResponse(ctx, http.StatusBadRequest, err)
		return
	}

	result := f.ServiceFactory().GenerateForecasts(solarSystemId, daysAmount)

	utils.SetResponse(ctx, http.StatusOK, result)
}

// ObtainForecast
// @Description Obtains a previously generated forecast
// @Success 201 {object} views.BaseResponse
// @Failure 400 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /forecast/solar_systems/{id}/obtain_forecasts [get]
func (f *ForecastController) ObtainForecast(ctx *gin.Context) {
	solarSystemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.SetResponse(ctx, http.StatusBadRequest, err)
		return
	}

	day, err := strconv.Atoi(ctx.Query("day"))
	if err != nil {
		utils.SetResponse(ctx, http.StatusBadRequest, err)
		return
	}
	result, forecastErr := f.ServiceFactory().ObtainForecast(solarSystemId, day)
	if forecastErr != nil {
		utils.SetResponse(ctx, utils.GetStatusErrorCode(forecastErr), forecastErr)
		return
	}

	utils.SetResponse(ctx, http.StatusOK, result)
}

// SoftDelete
// @Description Soft deletes existing systems, planets and forecasts
// @Success 201 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /forecast/all/soft [delete]
func (f *ForecastController) SoftDelete(ctx *gin.Context) {
	result, err := f.ServiceFactory().CleanData(true)
	if err != nil {
		utils.SetResponse(ctx, utils.GetStatusErrorCode(err), err)
		return
	}

	utils.SetResponse(ctx, http.StatusOK, result)
}

// SoftDelete
// @Description Soft deletes existing systems, planets and forecasts
// @Success 201 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /forecast/all/hard [delete]
func (f *ForecastController) HardDelete(ctx *gin.Context) {
	result, err := f.ServiceFactory().CleanData(false)
	if err != nil {
		utils.SetResponse(ctx, utils.GetStatusErrorCode(err), err)
		return
	}

	utils.SetResponse(ctx, http.StatusOK, result)
}
