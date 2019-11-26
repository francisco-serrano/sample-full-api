package models

import (
	"github.com/jinzhu/gorm"
	"math"
)

type SolarSystem struct {
	gorm.Model
	Name string `gorm:"column:name;not null"`
}

type Planet struct {
	gorm.Model
	Name          string  `gorm:"column:name;not null"`
	R             float64 `gorm:"column:r;not null"`
	Degrees       float64 `gorm:"column:degrees;not null"`
	Speed         float64 `gorm:"column:speed;not null"`
	Clockwise     bool    `gorm:"column:clockwise;not null"`
	Radians       float64 `gorm:"column:radians;not null"`
	X             float64 `gorm:"column:x;not null"`
	Y             float64 `gorm:"column:y;not null"`
	SolarSystemID uint    `gorm:"column:solar_system_id;not null"`
}

func (p *Planet) AdvanceDay() {
	if !p.Clockwise {
		newPosition := p.Degrees + p.Speed

		for newPosition >= 360 {
			newPosition -= 360
		}

		p.Degrees = newPosition
		p.Radians = newPosition * (math.Pi / 180.0)
		p.X = math.Round(p.R*math.Cos(p.Radians)*100) / 100
		p.X = math.Round(p.R*math.Sin(p.Radians)*100) / 100

		return
	}

	newPosition := p.Degrees - p.Speed

	for newPosition < 0 {
		newPosition += 360
	}

	p.Degrees = newPosition
	p.Radians = newPosition * (math.Pi / 180.0)
	p.X = math.Round(p.R*math.Cos(p.Radians)*100) / 100
	p.X = math.Round(p.R*math.Sin(p.Radians)*100) / 100
}

type DayForecast struct {
	gorm.Model
	SolarSystemID       uint    `gorm:"column:solar_system_id;not null"`
	Day                 int     `gorm:"column:day;not null"`
	RainIntensity       float64 `gorm:"column:rain_intensity;not null"`
	Drought             bool    `gorm:"column:drought;not null"`
	OptimalTempPressure bool    `gorm:"column:optimal_temp_pressure;not null"`
}
