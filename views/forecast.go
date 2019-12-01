package views

import "github.com/sample-full-api/models"

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

type BaseResponse struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

type GetForecastResponse struct {
	Day      int                `json:"day"`
	Forecast models.DayForecast `json:"forecast"`
}

type CleanDataResponse struct {
	Message string `json:"message"`
}
