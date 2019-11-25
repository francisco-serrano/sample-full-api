package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/services"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ForecastController struct {
	ServiceFactory func() *services.PlanetService
}

func (f *ForecastController) AddSolarSystem(ctx *gin.Context) {
	rawRequest, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}

	var request services.AddSolarSystemRequest
	if err = json.Unmarshal(rawRequest, &request); err != nil {
		panic(err)
	}

	result := f.ServiceFactory().AddSolarSystem(&request)

	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func (f *ForecastController) AddPlanet(ctx *gin.Context) {
	rawRequest, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}

	var request services.AddPlanetRequest
	if err = json.Unmarshal(rawRequest, &request); err != nil {
		panic(err)
	}

	result := f.ServiceFactory().AddPlanet(&request)

	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func (f *ForecastController) GenerateForecasts(ctx *gin.Context) {
	solarSystemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		panic(err)
	}

	result := f.ServiceFactory().GenerateForecasts(solarSystemId)

	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func (f *ForecastController) ObtainForecast(ctx *gin.Context) {
	day := ctx.Request.URL.Query().Get("day")

	if day == "" {
		panic("invalid day")
	}

	ctx.JSON(http.StatusOK, gin.H{
		"day":      day,
		"forecast": "rain",
	})
}
