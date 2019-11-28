package exercises

import (
	"fmt"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/utils"
	"strings"
)

type AnalysisResult struct {
	Droughts          int64
	RainyPeriods      int64
	MaxPeak           int64
	OptimalConditions int64
}

func AnalyzeDays(days int, logSituations bool, solarSystemID int, srcPlanets ...models.Planet) ([]interface{}, AnalysisResult) {
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

			if logSituations {
				fmt.Printf("drought detected at day %v\t\tpositions %s\n",
					day,
					strings.Join(positions, ";"),
				)
			}

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

			if logSituations {
				fmt.Printf("rainy period at day %v\t\t%s\n",
					day,
					strings.Join(positions, "\t"),
				)
			}

			amountRains += 1
			forecast.RainIntensity = perimeter
		}

		// exercise 3
		if utils.AlignedWithoutSun(planets...) && !utils.AlignedWithSun(planets...) {
			var positions []string
			for _, planet := range planets {
				positions = append(positions, fmt.Sprintf("r=%v, %v°", planet.R, planet.Degrees))
			}

			if logSituations {
				fmt.Printf("optimal condition detected at day %d with positions %s\n",
					day,
					strings.Join(positions, "\t"),
				)
			}

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
