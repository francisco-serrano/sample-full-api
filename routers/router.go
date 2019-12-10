package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sample-full-api/controllers"
	_ "github.com/sample-full-api/docs"
	"github.com/sample-full-api/middlewares"
	"github.com/sample-full-api/services"
	"github.com/sample-full-api/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	authController := controllers.AuthenticationController{
		AuthServiceFactory: func() controllers.AuthenticationService {
			return controllers.NewAuthenticationService(deps)
		},
	}

	forecastGroup := engine.Group("/forecast")
	forecastGroup.GET("/health", healthController.HealthCheck)
	forecastGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	forecastGroup.Use(middlewares.VerifyToken())
	forecastGroup.POST("/planets", planetController.AddPlanet)
	forecastGroup.GET("/planets", planetController.GetPlanets)
	forecastGroup.POST("/solar_systems", planetController.AddSolarSystem)
	forecastGroup.GET("/solar_systems", planetController.GetSolarSystems)
	forecastGroup.POST("/solar_systems/:id/generate_forecasts", planetController.GenerateForecasts)
	forecastGroup.GET("/solar_systems/:id/obtain_forecasts", planetController.ObtainForecast)
	forecastGroup.DELETE("/all/soft", planetController.SoftDelete)
	forecastGroup.DELETE("/all/hard", planetController.HardDelete)

	authGroup := engine.Group("/authentication")
	authGroup.POST("/login", authController.Login)
}
