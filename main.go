package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/routers"
	"github.com/sample-full-api/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func checkEnvironmentVariables() error {
	envVars := []string{"PORT", "LOG_LEVEL", "DB_USER", "DB_PASS", "DB_HOST"}

	for _, v := range envVars {
		if myVar := os.Getenv(v); myVar == "" {
			return errors.New(fmt.Sprintf("%s not provided", v))
		}
	}

	return nil
}

func obtainDbConnection() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")

	if user == "" || pass == "" || host == "" {
		return nil, errors.New("invalid user/pass/host")
	}

	connectionUrl := fmt.Sprintf("%s:%s@(%s)/solar_system_db?parseTime=true", user, pass, host)

	db, err := gorm.Open("mysql", connectionUrl)
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

func obtainLogger() (*log.Logger, error) {
	logger := log.New()
	logger.SetOutput(os.Stdout)

	switch os.Getenv("LOG_LEVEL") {
	case "INFO":
		logger.SetLevel(log.InfoLevel)
	case "ERROR":
		logger.SetLevel(log.ErrorLevel)
	case "DEBUG":
		logger.SetLevel(log.DebugLevel)
	default:
		return nil, errors.New("invalid log level value")
	}

	return logger, nil
}

func main() {
	if err := checkEnvironmentVariables(); err != nil {
		panic(err)
	}

	db, err := obtainDbConnection()
	if err != nil {
		panic(err)
	}

	logger, err := obtainLogger()
	if err != nil {
		panic(err)
	}

	deps := utils.Dependencies{
		Db:     db,
		Logger: logger,
	}

	engine := gin.Default()

	routers.InitializeRoutes(engine, deps)

	if err := engine.Run(); err != nil {
		panic(err)
	}
}
