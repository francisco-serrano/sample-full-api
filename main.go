package main

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/routers"
	"github.com/sample-full-api/services"
	"math"
	"time"
)

func obtainDbConnection() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root@/solar_system_db?parseTime=true")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.SolarSystem{}, &models.Planet{}, &models.DayForecast{})
	db.Model(&models.DayForecast{}).AddForeignKey("solar_system_id", "solar_systems(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Planet{}).AddForeignKey("solar_system_id", "solar_systems(id)", "RESTRICT", "RESTRICT")

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)

	return db
}

func main() {
	//ferengi, _ := models.NewPoint(500, 0, 1, true)
	//betasoide, _ := models.NewPoint(2000, 0, 3, true)
	//vulcano, _ := models.NewPoint(1000, 0, 5, false)
	//ferengi, _ := BuildPlanet(&services.AddPlanetRequest{
	//	Name:           "ferengi",
	//	Radio:          500,
	//	InitialDegrees: 120,
	//	SpeedByDay:     1,
	//	Clockwise:      true,
	//	SolarSystemId:  1,
	//})
	//
	//betasoide, _ := BuildPlanet(&services.AddPlanetRequest{
	//	Name:           "betasoide",
	//	Radio:          2000,
	//	InitialDegrees: 90,
	//	SpeedByDay:     3,
	//	Clockwise:      true,
	//	SolarSystemId:  1,
	//})
	//
	//vulcano, _ := BuildPlanet(&services.AddPlanetRequest{
	//	Name:           "vulcano",
	//	Radio:          1000,
	//	InitialDegrees: 30,
	//	SpeedByDay:     5,
	//	Clockwise:      false,
	//	SolarSystemId:  1,
	//})
	//
	//fmt.Printf("ferengi %+v\n", ferengi)
	//fmt.Printf("betasoide %+v\n", betasoide)
	//fmt.Printf("vulcano %+v\n", vulcano)
	//
	//planets := []models.Planet{*ferengi, *betasoide, *vulcano}
	//
	//aux1 := make([]models.Planet, len(planets))
	//aux2 := make([]models.Planet, len(planets))
	//aux3 := make([]models.Planet, len(planets))
	//copy(aux1, planets)
	//copy(aux2, planets)
	//copy(aux3, planets)
	//
	//exercises.AmountDroughts(365*10, false, aux1...)
	//exercises.AmountRainyPeriods(365 * 10, false, aux2...)
	//exercises.AmountOptimalConditions(365 * 10, false, aux3...)

	db := obtainDbConnection()
	router := routers.ObtainRoutes(db)

	if err := router.Run(); err != nil {
		panic(err)
	}
}

func BuildPlanet(req *services.AddPlanetRequest) (*models.Planet, error) {
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
