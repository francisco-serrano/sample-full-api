package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/services"
	"github.com/sample-full-api/utils"
	"github.com/sample-full-api/views"
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
// @Router /solar_systems [post]
func (f *ForecastController) AddSolarSystem(ctx *gin.Context) {
	var request views.AddSolarSystemRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
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

// GetSolarSystems
// @Description Gets all added solar systems
// @Success 201 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /solar_systems [get]
func (f *ForecastController) GetSolarSystems(ctx *gin.Context) {
	result, err := f.ServiceFactory().GetSolarSystems()
	if err != nil {
		utils.StatusInternalErrorResponse(ctx, err)
		return
	}

	utils.StatusOkResponse(ctx, result)
}

// AddPlanet
// @Description Adds a planet into the database
// @Param message body views.AddPlanetRequest false "Add Planet Request"
// @Success 201 {object} views.BaseResponse
// @Failure 400 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /planets [post]
func (f *ForecastController) AddPlanet(ctx *gin.Context) {
	var request views.AddPlanetRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
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

// GetPlanets
// @Description Gets all added planets
// @Success 201 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /planets [get]
func (f *ForecastController) GetPlanets(ctx *gin.Context) {
	result, err := f.ServiceFactory().GetPlanets()
	if err != nil {
		utils.StatusInternalErrorResponse(ctx, err)
		return
	}

	utils.StatusOkResponse(ctx, result)
}

// GenerateForecasts
// @Description Generates forecasts given added system and planets
// @Success 201 {object} views.BaseResponse
// @Failure 400 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /solar_systems/:id/generate_forecasts [post]
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

// ObtainForecast
// @Description Obtains a previously generated forecast
// @Success 201 {object} views.BaseResponse
// @Failure 400 {object} views.BaseResponse
// @Failure 500 {object} views.BaseResponse
// @Router /planets/forecast [get]
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

// SoftDelete
// @Description Soft deletes existing systems, planets and forecasts
// @Success 201 {object} views.BaseResponse
// @Router /all/soft [delete]
func (f *ForecastController) SoftDelete(ctx *gin.Context) {
	result, err := f.ServiceFactory().CleanData(true)
	if err != nil {
		utils.StatusInternalErrorResponse(ctx, err)
		return
	}

	utils.StatusOkResponse(ctx, result)
}

// SoftDelete
// @Description Soft deletes existing systems, planets and forecasts
// @Success 201 {object} views.BaseResponse
// @Router /all/hard [delete]
func (f *ForecastController) HardDelete(ctx *gin.Context) {
	result, err := f.ServiceFactory().CleanData(false)
	if err != nil {
		utils.StatusInternalErrorResponse(ctx, err)
		return
	}

	utils.StatusOkResponse(ctx, result)
}
