package views

import "github.com/sample-full-api/models"

type AddPlanetRequest struct {
	Name            string  `json:"name"`
	Radio           float64 `json:"radio"`
	InitialDegrees  float64 `json:"initial_degrees"`
	SpeedByDay      float64 `json:"speed_by_day"`
	Clockwise       bool    `json:"clockwise"`
	SolarSystemName string  `json:"solar_system_name"`
}

type AddPlanetResponse struct {
	Name            string  `json:"name"`
	Radio           float64 `json:"radio"`
	InitialDegrees  float64 `json:"initial_degrees"`
	SpeedByDay      float64 `json:"speed_by_day"`
	Clockwise       bool    `json:"clockwise"`
	SolarSystemName string  `json:"solar_system_name"`
	X               float64 `json:"x"`
	Y               float64 `json:"y"`
}

type GetPlanetResponse struct {
	Name           string  `json:"name"`
	Radio          float64 `json:"radio"`
	InitialDegrees float64 `json:"initial_degrees"`
	SpeedByDay     float64 `json:"speed_by_day"`
	Clockwise      bool    `json:"clockwise"`
	X              float64 `json:"x"`
	Y              float64 `json:"y"`
	SolarSystemID  uint    `json:"solar_system_name"`
}

type GetPlanetsResponse struct {
	Planets []GetPlanetResponse `json:"planets"`
}

type AddSolarSystemRequest struct {
	Name string `json:"name"`
}

type AddSolarSystemResponse struct {
	Name string `json:"name"`
}

type GetSolarSystemResponse struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
}

type GetSolarSystemsResponse struct {
	SolarSystems []GetSolarSystemResponse `json:"solar_systems"`
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
