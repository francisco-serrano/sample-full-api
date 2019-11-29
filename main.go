package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/routers"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func checkEnvironmentVariables() {
	envVars := []string{"PORT"}

	for _, v := range envVars {
		if myVar := os.Getenv(v); myVar == "" {
			panic(fmt.Sprintf("%s not provided", v))
		}
	}
}

func obtainDbConnection() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:root@/solar_system_db?parseTime=true")
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.SolarSystem{}, &models.Planet{}, &models.DayForecast{})
	db.Model(&models.DayForecast{}).AddForeignKey("solar_system_id", "solar_systems(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Planet{}).AddForeignKey("solar_system_id", "solar_systems(id)", "RESTRICT", "RESTRICT")

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)

	return db, nil
}

func main() {
	checkEnvironmentVariables()

	db, err := obtainDbConnection()
	if err != nil {
		panic(err)
	}

	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel)

	router := routers.ObtainRoutes(db, logger)

	if err := router.Run(); err != nil {
		panic(err)
	}
}
