package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/controllers"
	"github.com/sample-full-api/services"
	"github.com/sample-full-api/utils"
)

func ObtainRoutes(deps utils.Dependencies) *gin.Engine {
	router := gin.Default()

	planetController := controllers.ForecastController{
		ServiceFactory: func() services.PlanetService {
			return services.NewPlanetService(deps.Db, deps.Logger)
		},
	}

	router.POST("/planets", planetController.AddPlanet)
	router.GET("/planets", planetController.GetPlanets)
	router.GET("/planets/forecast", planetController.ObtainForecast)

	router.POST("/solar_systems", planetController.AddSolarSystem)
	router.GET("/solar_systems", planetController.GetSolarSystems)
	router.POST("/solar_systems/:id/generate_forecasts", planetController.GenerateForecasts)

	return router
}
