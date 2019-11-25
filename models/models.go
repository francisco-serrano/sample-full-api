package models

import "github.com/jinzhu/gorm"

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

type DayForecast struct {
	gorm.Model
	PlanetID            uint `gorm:"column:planet_id;not null"`
	Day                 int  `gorm:"column:day;not null"`
	RainIntensity       int  `gorm:"column:rain_intensity;not null"`
	Drought             bool `gorm:"column:drought;not null"`
	OptimalTempPressure bool `gorm:"column:optimal_temp_pressure;not null"`
}
