package services

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/exercises"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/views"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
	"math"
)

type PlanetService interface {
	AddPlanet(request *views.AddPlanetRequest) *models.Planet
	AddSolarSystem(request *views.AddSolarSystemRequest) *models.SolarSystem
	GenerateForecasts(solarSystemId, daysAmount int) string
	ObtainForecast(day int) gin.H
}

type planetService struct {
	db *gorm.DB
}

func NewPlanetService(db *gorm.DB) *planetService {
	return &planetService{db: db}
}

func (p *planetService) AddPlanet(request *views.AddPlanetRequest) *models.Planet {
	planet, err := p.buildPlanet(request)
	if err != nil {
		panic(err)
	}

	if err = p.db.Create(planet).Error; err != nil {
		panic(err)
	}

	return planet
}

func (p *planetService) AddSolarSystem(request *views.AddSolarSystemRequest) *models.SolarSystem {
	solarSystem, err := p.buildSolarSystem(request)
	if err != nil {
		panic(err)
	}

	if err = p.db.Create(solarSystem).Error; err != nil {
		panic(err)
	}

	return solarSystem
}

func (p *planetService) GenerateForecasts(solarSystemId, daysAmount int) string {
	go p.generateForecast(solarSystemId, daysAmount)

	return fmt.Sprintf("job triggered for system %d", solarSystemId)
}

func (p *planetService) ObtainForecast(day int) gin.H {
	var forecast models.DayForecast
	forecast.Day = day
	forecast.DeletedAt = nil

	if err := p.db.First(&forecast).Error; err != nil {
		panic(err)
	}

	return gin.H{
		"day":      day,
		"forecast": forecast,
	}
}

func (p *planetService) buildSolarSystem(req *views.AddSolarSystemRequest) (*models.SolarSystem, error) {
	return &models.SolarSystem{
		Name: req.Name,
	}, nil
}

func (p *planetService) buildPlanet(req *views.AddPlanetRequest) (*models.Planet, error) {
	if req.Radio < 0.0 || req.InitialDegrees < 0.0 || req.InitialDegrees >= 360.0 {
		return nil, errors.New("invalid input data")
	}

	radians := req.InitialDegrees * (math.Pi / 180.0)

	return &models.Planet{
		Name:          req.Name,
		R:             req.Radio,
		Degrees:       req.InitialDegrees,
		Speed:         req.SpeedByDay,
		Clockwise:     req.Clockwise,
		Radians:       radians,
		X:             math.Round(req.Radio*math.Cos(radians)*100) / 100,
		Y:             math.Round(req.Radio*math.Sin(radians)*100) / 100,
		SolarSystemID: req.SolarSystemId,
	}, nil
}

// goroutine
func (p *planetService) generateForecast(solarSystemId, daysAmount int) {
	fmt.Printf("generating forecast for system %d for %d days\n", solarSystemId, daysAmount)

	p.cleanUpExistingForecasts(solarSystemId)

	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		panic(err)
	}

	var planets []models.Planet
	if err := tx.Where(&models.Planet{SolarSystemID: uint(solarSystemId)}).Find(&planets).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	planetsCopy := make([]models.Planet, len(planets))
	copy(planetsCopy, planets)

	forecasts, result := exercises.AnalyzeDays(daysAmount, false, solarSystemId, planetsCopy...)

	fmt.Printf("result %+v\n", result)

	if err := gormbulk.BulkInsert(tx, forecasts, 2000); err != nil {
		tx.Rollback()
		panic(err)
	}

	if err := tx.Commit().Error; err != nil {
		panic(err)
	}
}

func (p *planetService) cleanUpExistingForecasts(solarSystemId int) {
	var existingForecasts []models.DayForecast
	if err := p.db.Find(&existingForecasts, "solar_system_id = ?", solarSystemId).Error; err != nil {
		panic(err)
	}

	if err := p.db.Delete(&existingForecasts).Error; err != nil {
		panic(err)
	}
}
