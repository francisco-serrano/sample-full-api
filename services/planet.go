package services

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/utils"
	"github.com/sample-full-api/views"
	log "github.com/sirupsen/logrus"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
	"math"
	"strings"
)

type PlanetService interface {
	AddPlanet(request *views.AddPlanetRequest) *models.Planet
	GetPlanets() []models.Planet
	AddSolarSystem(request *views.AddSolarSystemRequest) *models.SolarSystem
	GetSolarSystems() []models.SolarSystem
	GenerateForecasts(solarSystemId, daysAmount int) string
	ObtainForecast(day int) gin.H
}

type planetService struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewPlanetService(db *gorm.DB, logger *log.Logger) *planetService {
	return &planetService{
		db:     db,
		logger: logger,
	}
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

func (p *planetService) GetPlanets() []models.Planet {
	var planets []models.Planet
	if err := p.db.Find(&planets).Error; err != nil {
		panic(err)
	}

	return planets
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

func (p *planetService) GetSolarSystems() []models.SolarSystem {
	var systems []models.SolarSystem
	if err := p.db.Find(&systems).Error; err != nil {
		panic(err)
	}

	return systems
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
	p.logger.Infof("generating forecast for system %d for %d days\n", solarSystemId, daysAmount)

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

	forecasts, result := p.analyzeDays(daysAmount, solarSystemId, planetsCopy...)

	p.logger.Infof("result %+v\n", result)

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

type AnalysisResult struct {
	Droughts          int64
	RainyPeriods      int64
	MaxPeak           int64
	OptimalConditions int64
}

func (p *planetService) analyzeDays(days int, solarSystemID int, srcPlanets ...models.Planet) ([]interface{}, AnalysisResult) {
	var forecasts []interface{}

	planets := make([]models.Planet, len(srcPlanets))
	copy(planets, srcPlanets)

	sun := models.Planet{
		Name:      "sun",
		R:         0,
		Degrees:   0,
		Speed:     0,
		Clockwise: false,
		Radians:   0,
		X:         0,
		Y:         0,
	}

	amountAlignments := 0
	amountRains, maxPerimeter, maxPerimeterDay := 0, 0.0, 0
	amount := 0
	for day := 0; day < days; day++ {
		var forecast models.DayForecast
		forecast.Day = day
		forecast.SolarSystemID = uint(solarSystemID)

		// exercise 1
		if utils.AlignedWithSun(planets...) {
			var positions []string
			for _, planet := range planets {
				positions = append(positions, fmt.Sprintf("%v", planet.Degrees))
			}

			p.logger.Debugf("drought detected at day %v\t\tpositions %s\n", day, strings.Join(positions, ";"))

			amountAlignments += 1
			forecast.Drought = true
		}

		// exercise 2
		if utils.WithinPolygon(sun, planets...) {
			var positions []string
			for _, planet := range planets {
				positions = append(positions, fmt.Sprintf("r=%v, %v°", planet.R, planet.Degrees))
			}

			perimeter := utils.Perimeter(planets...)

			if perimeter > maxPerimeter {
				maxPerimeter = perimeter
				maxPerimeterDay = day
			}

			p.logger.Debugf("rainy period at day %v\t\t%s\n", day, strings.Join(positions, "\t"))

			amountRains += 1
			forecast.RainIntensity = perimeter
		}

		// exercise 3
		if utils.AlignedWithoutSun(planets...) && !utils.AlignedWithSun(planets...) {
			var positions []string
			for _, planet := range planets {
				positions = append(positions, fmt.Sprintf("r=%v, %v°", planet.R, planet.Degrees))
			}

			p.logger.Debugf("optimal condition detected at day %d with positions %s\n", day, strings.Join(positions, "\t"))

			amount += 1
			forecast.OptimalTempPressure = true
		}

		forecasts = append(forecasts, forecast)

		for i := 0; i < len(planets); i++ {
			planets[i].AdvanceDay()
		}
	}

	result := AnalysisResult{
		Droughts:          int64(amountAlignments),
		RainyPeriods:      int64(amountRains),
		MaxPeak:           int64(maxPerimeterDay),
		OptimalConditions: int64(amount),
	}

	return forecasts, result
}
