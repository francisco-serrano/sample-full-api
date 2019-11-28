package services

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/exercises"
	"github.com/sample-full-api/models"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
	"math"
	"time"
)

type PlanetService struct {
	db *gorm.DB
}

func NewPlanetService(db *gorm.DB) *PlanetService {
	return &PlanetService{db: db}
}

func (p *PlanetService) AddPlanet(request *AddPlanetRequest) *models.Planet {
	planet, err := p.buildPlanet(request)
	if err != nil {
		panic(err)
	}

	if err = p.db.Create(planet).Error; err != nil {
		panic(err)
	}

	return planet
}

func (p *PlanetService) AddSolarSystem(request *AddSolarSystemRequest) *models.SolarSystem {
	solarSystem, err := p.buildSolarSystem(request)
	if err != nil {
		panic(err)
	}

	if err = p.db.Create(solarSystem).Error; err != nil {
		panic(err)
	}

	return solarSystem
}

func (p *PlanetService) GenerateForecasts(solarSystemId, daysAmount int) string {
	go p.generateForecast(solarSystemId, daysAmount)

	return fmt.Sprintf("job triggered for system %d", solarSystemId)
}

func (p *PlanetService) ObtainForecast(day int) gin.H {
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

func (p *PlanetService) buildSolarSystem(req *AddSolarSystemRequest) (*models.SolarSystem, error) {
	return &models.SolarSystem{
		Name: req.Name,
	}, nil
}

func (p *PlanetService) buildPlanet(req *AddPlanetRequest) (*models.Planet, error) {
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
func (p *PlanetService) generateForecast(solarSystemId, daysAmount int) {
	time.Sleep(1 * time.Second)

	if solarSystemId < 0 {
		panic("solar system cannot be negative")
	}

	fmt.Printf("generating forecast for system %d for %d days\n", solarSystemId, daysAmount)

	var planets []models.Planet
	if err := p.db.Where(&models.Planet{SolarSystemID: uint(solarSystemId)}).Find(&planets).Error; err != nil {
		panic(err)
	}

	aux1 := make([]models.Planet, len(planets))
	copy(aux1, planets)

	forecasts, result := exercises.AnalyzeDays(daysAmount, false, solarSystemId, aux1...)

	fmt.Printf("analysis result %+v\n", result)

	var existingForecasts []models.DayForecast
	if err := p.db.Find(&existingForecasts, "solar_system_id = ?", solarSystemId).Error; err != nil {
		panic(err)
	}

	fmt.Println("AAAAAAAAAAAAAAA", len(existingForecasts))

	if err := p.db.Delete(&existingForecasts).Error; err != nil {
		panic(err)
	}

	if err := gormbulk.BulkInsert(p.db, forecasts, 2000); err != nil {
		panic(err)
	}
}

// ----- VIEWS -----
type AddPlanetRequest struct {
	Name           string  `json:"name"`
	Radio          float64 `json:"radio"`
	InitialDegrees float64 `json:"initial_degrees"`
	SpeedByDay     float64 `json:"speed_by_day"`
	Clockwise      bool    `json:"clockwise"`
	SolarSystemId  uint    `json:"solar_system_id"`
}

type AddSolarSystemRequest struct {
	Name string `json:"name"`
}
