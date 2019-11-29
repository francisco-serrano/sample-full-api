package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/controllers"
	"github.com/sample-full-api/services"
)

func ObtainRoutes(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	planetController := controllers.ForecastController{
		ServiceFactory: func() services.PlanetService {
			return services.NewPlanetService(db)
		},
	}

	router.POST("/planets", planetController.AddPlanet)
	router.POST("/solar_system", planetController.AddSolarSystem)
	router.POST("/solar_system/:id/generate_forecasts", planetController.GenerateForecasts)
	router.GET("/planets/forecast", planetController.ObtainForecast)

	return router
}
