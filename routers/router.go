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
	planetController := controllers.ForecastController{
		ServiceFactory: func() services.PlanetService {
			return services.NewPlanetService(deps.Db, deps.Logger)
		},
	}

	engine.POST("/planets", planetController.AddPlanet)
	engine.GET("/planets", planetController.GetPlanets)
	engine.GET("/planets/forecast", planetController.ObtainForecast)

	engine.POST("/solar_systems", planetController.AddSolarSystem)
	engine.GET("/solar_systems", planetController.GetSolarSystems)
	engine.POST("/solar_systems/:id/generate_forecasts", planetController.GenerateForecasts)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
