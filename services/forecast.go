package services

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/utils"
	"github.com/sample-full-api/views"
	log "github.com/sirupsen/logrus"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
	"math"
	"strings"
)

type ForecastService interface {
	AddPlanet(request *views.AddPlanetRequest) (*models.Planet, error)
	GetPlanets() (*[]models.Planet, error)
	AddSolarSystem(request *views.AddSolarSystemRequest) (*models.SolarSystem, error)
	GetSolarSystems() (*[]models.SolarSystem, error)
	GenerateForecasts(solarSystemId, daysAmount int) string
	ObtainForecast(day int) (*views.GetForecastResponse, error)
	CleanData(softDelete bool) (*views.CleanDataResponse, error)
}

type forecastService struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewPlanetService(deps utils.Dependencies) *forecastService {
	return &forecastService{
		db:     deps.Db,
		logger: deps.Logger,
	}
}

func (p *forecastService) AddPlanet(request *views.AddPlanetRequest) (*models.Planet, error) {
	planet, err := p.buildPlanet(request)
	if err != nil {
		return nil, err
	}

	if err = p.db.Create(planet).Error; err != nil {
		return nil, err
	}

	return planet, nil
}

func (p *forecastService) GetPlanets() (*[]models.Planet, error) {
	var planets []models.Planet
	if err := p.db.Find(&planets).Error; err != nil {
		return nil, err
	}

	return &planets, nil
}

func (p *forecastService) AddSolarSystem(request *views.AddSolarSystemRequest) (*models.SolarSystem, error) {
	solarSystem, err := p.buildSolarSystem(request)
	if err != nil {
		return nil, err
	}

	if err = p.db.Create(solarSystem).Error; err != nil {
		return nil, err
	}

	return solarSystem, nil
}

func (p *forecastService) GetSolarSystems() (*[]models.SolarSystem, error) {
	var systems []models.SolarSystem
	if err := p.db.Find(&systems).Error; err != nil {
		return nil, err
	}

	return &systems, nil
}

func (p *forecastService) GenerateForecasts(solarSystemId, daysAmount int) string {
	go p.generateForecast(solarSystemId, daysAmount)

	return fmt.Sprintf("job triggered for system %d", solarSystemId)
}

func (p *forecastService) ObtainForecast(day int) (*views.GetForecastResponse, error) {
	var forecast models.DayForecast
	forecast.Day = day
	forecast.DeletedAt = nil

	if err := p.db.First(&forecast).Error; err != nil {
		return nil, err
	}

	return &views.GetForecastResponse{
		Day:      day,
		Forecast: forecast,
	}, nil
}

func (p *forecastService) CleanData(softDelete bool) (*views.CleanDataResponse, error) {
	if softDelete {
		if err := p.transactionalSoftDelete(); err != nil {
			return nil, err
		}

		return &views.CleanDataResponse{Message: "soft delete performed successfully"}, nil
	}

	if err := p.transactionalHardDelete(); err != nil {
		return nil, err
	}

	return &views.CleanDataResponse{Message: "hard delete performed successfully"}, nil
}

func (p *forecastService) transactionalSoftDelete() error {
	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := p.db.Delete(&models.DayForecast{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := p.db.Delete(&models.Planet{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := p.db.Delete(&models.SolarSystem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (p *forecastService) transactionalHardDelete() error {
	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := p.db.Unscoped().Delete(&models.DayForecast{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := p.db.Unscoped().Delete(&models.Planet{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := p.db.Unscoped().Delete(&models.SolarSystem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (p *forecastService) buildSolarSystem(req *views.AddSolarSystemRequest) (*models.SolarSystem, error) {
	return &models.SolarSystem{
		Name: req.Name,
	}, nil
}

func (p *forecastService) buildPlanet(req *views.AddPlanetRequest) (*models.Planet, error) {
	var solarSystem models.SolarSystem
	solarSystem.Name = req.SolarSystemName
	if err := p.db.Find(&solarSystem).Error; err != nil {
		return nil, err
	}

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
		SolarSystemID: solarSystem.ID,
	}, nil
}

// goroutine
func (p *forecastService) generateForecast(solarSystemId, daysAmount int) error {
	p.logger.Infof("generating forecast for system %d for %d days\n", solarSystemId, daysAmount)

	if err := p.cleanUpExistingForecasts(solarSystemId); err != nil {
		return err
	}

	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	var planets []models.Planet
	if err := tx.Where(&models.Planet{SolarSystemID: uint(solarSystemId)}).Find(&planets).Error; err != nil {
		tx.Rollback()
	}

	planetsCopy := make([]models.Planet, len(planets))
	copy(planetsCopy, planets)

	forecasts, result := p.analyzeDays(daysAmount, solarSystemId, planetsCopy...)

	p.logger.Infof("result %+v\n", result)

	if err := gormbulk.BulkInsert(tx, forecasts, 2000); err != nil {
		tx.Rollback()
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (p *forecastService) cleanUpExistingForecasts(solarSystemId int) error {
	var existingForecasts []models.DayForecast
	if err := p.db.Find(&existingForecasts, "solar_system_id = ?", solarSystemId).Error; err != nil {
		return err
	}

	if err := p.db.Delete(&existingForecasts).Error; err != nil {
		return err
	}

	return nil
}

type AnalysisResult struct {
	Droughts          int64
	RainyPeriods      int64
	MaxPeak           int64
	OptimalConditions int64
}

func (p *forecastService) analyzeDays(days int, solarSystemID int, srcPlanets ...models.Planet) ([]interface{}, AnalysisResult) {
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
