package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/controllers"
	"github.com/sample-full-api/services"
	"github.com/sample-full-api/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/sample-full-api/docs"
)

func InitializeRoutes(engine *gin.Engine, deps utils.Dependencies) {
	healthController := controllers.HealthController{
		ServiceFactory: func() services.HealthService {
			return services.NewHealthService(deps)
		},
	}

	planetController := controllers.ForecastController{
		ServiceFactory: func() services.ForecastService {
			return services.NewPlanetService(deps)
		},
	}

	group := engine.Group("/forecast")

	group.GET("/health", healthController.HealthCheck)

	group.POST("/planets", planetController.AddPlanet)
	group.GET("/planets", planetController.GetPlanets)

	group.POST("/solar_systems", planetController.AddSolarSystem)
	group.GET("/solar_systems", planetController.GetSolarSystems)

	group.POST("/solar_systems/:id/generate_forecasts", planetController.GenerateForecasts)
	group.GET("/solar_systems/:id/obtain_forecasts", planetController.ObtainForecast)

	group.DELETE("/all/soft", planetController.SoftDelete)
	group.DELETE("/all/hard", planetController.HardDelete)

	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
