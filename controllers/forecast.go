package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type ForecastController struct {
	ServiceFactory func() *PlanetService
}

func (f *ForecastController) AddPlanet(ctx *gin.Context) {
	f.ServiceFactory().AddPlanet()

	ctx.JSON(http.StatusOK, gin.H{
		"msg": `planet added`,
	})
}

func (f *ForecastController) GenerateForecasts(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "job triggered",
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

type PlanetService struct {
	Db *gorm.DB
}

func NewPlanetService(db *gorm.DB) *PlanetService {
	return &PlanetService{Db: db}
}

func (p *PlanetService) AddPlanet() {
	fmt.Println("adding planet...")
}
