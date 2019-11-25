package main

import (
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/routers"

	_ "github.com/go-sql-driver/mysql"
)

func obtainDbConnection() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root@/solar_system_db?parseTime=true")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.SolarSystem{}, &models.Planet{}, &models.DayForecast{})
	db.Model(&models.DayForecast{}).AddForeignKey("planet_id", "planets(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Planet{}).AddForeignKey("solar_system_id", "solar_systems(id)", "RESTRICT", "RESTRICT")

	return db
}

func main() {
	//exercises.AmountDroughts(365*10, false)
	//exercises.AmountRainyPeriods(365*10, false)
	//exercises.AmountOptimalConditions(365*10, false)

	db := obtainDbConnection()
	router := routers.ObtainRoutes(db)

	if err := router.Run(); err != nil {
		panic(err)
	}
}
