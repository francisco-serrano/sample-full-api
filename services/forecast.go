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
	AddPlanet(request *views.AddPlanetRequest) (*views.AddPlanetResponse, *utils.ForecastError)
	GetPlanets() (*views.GetPlanetsResponse, *utils.ForecastError)
	AddSolarSystem(request *views.AddSolarSystemRequest) (*views.AddSolarSystemResponse, *utils.ForecastError)
	GetSolarSystems() (*views.GetSolarSystemsResponse, *utils.ForecastError)
	GenerateForecasts(solarSystemId, daysAmount int) string
	ObtainForecast(solarSystemId, day int) (*views.GetForecastResponse, *utils.ForecastError)
	CleanData(softDelete bool) (*views.CleanDataResponse, *utils.ForecastError)
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

func (p *forecastService) AddPlanet(request *views.AddPlanetRequest) (*views.AddPlanetResponse, *utils.ForecastError) {
	planet, err := p.buildPlanet(request)
	if err != nil {
		p.logger.Error(err)
		return nil, err
	}

	if err := p.db.Create(planet).Error; err != nil {
		p.logger.Error(err)
		return nil, utils.ErrorInternal(err.Error())
	}

	return &views.AddPlanetResponse{
		Name:            planet.Name,
		Radio:           planet.R,
		InitialDegrees:  planet.Degrees,
		SpeedByDay:      planet.Speed,
		Clockwise:       planet.Clockwise,
		SolarSystemName: request.SolarSystemName,
		X:               planet.X,
		Y:               planet.Y,
	}, nil
}

func (p *forecastService) GetPlanets() (*views.GetPlanetsResponse, *utils.ForecastError) {
	var planets []models.Planet
	if err := p.db.Find(&planets).Error; err != nil {
		p.logger.Error(err)
		return nil, utils.ErrorInternal(err.Error())
	}

	var planetsResponse []views.GetPlanetResponse
	for _, planet := range planets {
		planetResponse := views.GetPlanetResponse{
			Name:           planet.Name,
			Radio:          planet.R,
			InitialDegrees: planet.Degrees,
			SpeedByDay:     planet.Speed,
			Clockwise:      planet.Clockwise,
			X:              planet.X,
			Y:              planet.Y,
			SolarSystemID:  planet.SolarSystemID,
		}
		planetsResponse = append(planetsResponse, planetResponse)
	}

	return &views.GetPlanetsResponse{Planets: planetsResponse}, nil
}

func (p *forecastService) AddSolarSystem(request *views.AddSolarSystemRequest) (*views.AddSolarSystemResponse, *utils.ForecastError) {
	solarSystem, err := p.buildSolarSystem(request)
	if err != nil {
		p.logger.Error(err)
		return nil, err
	}

	if err := p.db.Create(solarSystem).Error; err != nil {
		p.logger.Error(err)
		return nil, utils.ErrorInternal("unable to add solar system")
	}

	return &views.AddSolarSystemResponse{Name: solarSystem.Name}, nil
}

func (p *forecastService) GetSolarSystems() (*views.GetSolarSystemsResponse, *utils.ForecastError) {
	var systems []models.SolarSystem
	if err := p.db.Find(&systems).Error; err != nil {
		p.logger.Error(err)
		return nil, utils.ErrorInternal(err.Error())
	}

	var systemsResponse []views.GetSolarSystemResponse
	for _, system := range systems {
		systemResponse := views.GetSolarSystemResponse{
			Name: system.Name,
			ID:   system.ID,
		}
		systemsResponse = append(systemsResponse, systemResponse)
	}

	return &views.GetSolarSystemsResponse{SolarSystems: systemsResponse}, nil
}

func (p *forecastService) GenerateForecasts(solarSystemId, daysAmount int) string {
	go func() {
		if err := p.generateForecast(solarSystemId, daysAmount); err != nil {
			p.logger.Error(err)
		}
	}()

	return fmt.Sprintf("job triggered for system %d", solarSystemId)
}

func (p *forecastService) ObtainForecast(solarSystemId, day int) (*views.GetForecastResponse, *utils.ForecastError) {
	var forecast models.DayForecast
	if err := p.db.Where(&models.DayForecast{
		Day:           day,
		SolarSystemID: uint(solarSystemId),
	}).First(&forecast).Error; err != nil {
		p.logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrorNotFound(err.Error())
		}

		return nil, utils.ErrorInternal(err.Error())
	}

	return &views.GetForecastResponse{
		Day:      day,
		Forecast: forecast,
	}, nil
}

func (p *forecastService) CleanData(softDelete bool) (*views.CleanDataResponse, *utils.ForecastError) {
	if softDelete {
		if err := p.transactionalSoftDelete(); err != nil {
			p.logger.Error(err)
			return nil, err
		}

		return &views.CleanDataResponse{Message: "soft delete performed successfully"}, nil
	}

	if err := p.transactionalHardDelete(); err != nil {
		p.logger.Error(err)
		return nil, err
	}

	return &views.CleanDataResponse{Message: "hard delete performed successfully"}, nil
}

func (p *forecastService) transactionalSoftDelete() *utils.ForecastError {
	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		p.logger.Error(err)
		return utils.ErrorInternal(err.Error())
	}

	if err := p.db.Delete(&models.DayForecast{}).Error; err != nil {
		p.logger.Error(err)
		tx.Rollback()
		return utils.ErrorInternal(err.Error())
	}

	if err := p.db.Delete(&models.Planet{}).Error; err != nil {
		p.logger.Error(err)
		tx.Rollback()
		return utils.ErrorInternal(err.Error())
	}

	if err := p.db.Delete(&models.SolarSystem{}).Error; err != nil {
		p.logger.Error(err)
		tx.Rollback()
		return utils.ErrorInternal(err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		p.logger.Error(err)
		return utils.ErrorInternal(err.Error())
	}

	return nil
}

func (p *forecastService) transactionalHardDelete() *utils.ForecastError {
	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		p.logger.Error(err)
		return utils.ErrorInternal(err.Error())
	}

	if err := p.db.Unscoped().Delete(&models.DayForecast{}).Error; err != nil {
		p.logger.Error(err)
		tx.Rollback()
		return utils.ErrorInternal(err.Error())
	}

	if err := p.db.Unscoped().Delete(&models.Planet{}).Error; err != nil {
		p.logger.Error(err)
		tx.Rollback()
		return utils.ErrorInternal(err.Error())
	}

	if err := p.db.Unscoped().Delete(&models.SolarSystem{}).Error; err != nil {
		p.logger.Error(err)
		tx.Rollback()
		return utils.ErrorInternal(err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		p.logger.Error(err)
		return utils.ErrorInternal(err.Error())
	}

	return nil
}

func (p *forecastService) buildSolarSystem(req *views.AddSolarSystemRequest) (*models.SolarSystem, *utils.ForecastError) {
	return &models.SolarSystem{
		Name: req.Name,
	}, nil
}

func (p *forecastService) buildPlanet(req *views.AddPlanetRequest) (*models.Planet, *utils.ForecastError) {
	var solarSystem models.SolarSystem
	if err := p.db.Where(&models.SolarSystem{
		Name: req.SolarSystemName,
	}).Find(&solarSystem).Error; err != nil {
		p.logger.Error(err)
		return nil, utils.ErrorInternal(err.Error())
	}

	if req.Radio < 0.0 || req.InitialDegrees < 0.0 || req.InitialDegrees >= 360.0 {
		err := errors.New("invalid input data")
		p.logger.Error(err)
		return nil, utils.ErrorInternal(err.Error())
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
func (p *forecastService) generateForecast(solarSystemId, daysAmount int) *utils.ForecastError {
	p.logger.Infof("generating forecast for system %d for %d days\n", solarSystemId, daysAmount)

	if err := p.cleanUpExistingForecasts(solarSystemId); err != nil {
		p.logger.Error(err)
		return err
	}

	tx := p.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		p.logger.Error(err)
		return utils.ErrorInternal(err.Error())
	}

	var planets []models.Planet
	if err := tx.Where(&models.Planet{SolarSystemID: uint(solarSystemId)}).Find(&planets).Error; err != nil {
		p.logger.Error(err)
		tx.Rollback()
	}

	planetsCopy := make([]models.Planet, len(planets))
	copy(planetsCopy, planets)

	forecasts, result := p.analyzeDays(daysAmount, solarSystemId, planetsCopy...)

	p.logger.Infof("result %+v\n", result)

	if err := gormbulk.BulkInsert(tx, forecasts, 2000); err != nil {
		p.logger.Error(err)
		tx.Rollback()
	}

	if err := tx.Commit().Error; err != nil {
		p.logger.Error(err)
		return utils.ErrorInternal(err.Error())
	}

	return nil
}

func (p *forecastService) cleanUpExistingForecasts(solarSystemId int) *utils.ForecastError {
	if err := p.db.Where(&models.DayForecast{
		SolarSystemID: uint(solarSystemId),
	}).Delete(&models.DayForecast{}).Error; err != nil {
		p.logger.Error(err)
		return utils.ErrorInternal(err.Error())
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
