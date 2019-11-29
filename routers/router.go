package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/controllers"
	"github.com/sample-full-api/services"
	log "github.com/sirupsen/logrus"
)

func ObtainRoutes(db *gorm.DB, logger *log.Logger) *gin.Engine {
	router := gin.Default()

	planetController := controllers.ForecastController{
		ServiceFactory: func() services.PlanetService {
			return services.NewPlanetService(db, logger)
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
