package services

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/exercises"
	"github.com/sample-full-api/models"
	"math"
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

func (p *PlanetService) GenerateForecasts(solarSystemId int) string {
	go p.generateForecast(solarSystemId)

	return fmt.Sprintf("job triggered for system %d", solarSystemId)
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

func (p *PlanetService) buildSolarSystem(req *AddSolarSystemRequest) (*models.SolarSystem, error) {
	return &models.SolarSystem{
		Name: req.Name,
	}, nil
}

// goroutine
func (p *PlanetService) generateForecast(solarSystemId int) {
	exercises.AmountDroughts(365*10, false)
	exercises.AmountRainyPeriods(365*10, false)
	exercises.AmountOptimalConditions(365*10, false)

	if solarSystemId < 0 {
		panic("solar system cannot be negative")
	}

	var planets []models.Planet
	if err := p.db.Where(&models.Planet{SolarSystemID: uint(solarSystemId)}).Find(&planets).Error; err != nil {
		panic(err)
	}

	for _, planet := range planets {
		fmt.Printf("%+v\n", planet)
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
