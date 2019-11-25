package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/controllers"
)

func ObtainRoutes(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	planetController := controllers.ForecastController{
		ServiceFactory: func() *controllers.PlanetService {
			return controllers.NewPlanetService(db)
		},
	}

	router.POST("/planets", planetController.AddPlanet)
	router.POST("/planets/generate_forecasts", planetController.GenerateForecasts)
	router.GET("/planets/forecast", planetController.ObtainForecast)

	return router
}
