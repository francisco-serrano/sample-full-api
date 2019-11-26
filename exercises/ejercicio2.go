package exercises

import (
	"fmt"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/utils"
	"strings"
)

func AmountRainyPeriods(days int, logSituations bool, srcPlanets ...models.Planet) {
	//ferengi, _ := models.NewPoint(500, 45, 1, true)
	//betasoide, _ := models.NewPoint(2000, 270, 3, true)
	//vulcano, _ := models.NewPoint(1000, 135, 5, false)
	//
	//planets := []*models.Point{ferengi, betasoide, vulcano}

	//sun, _ := models.NewPoint(0, 0, 0, false)
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

	amountRains, maxPerimeter, maxPerimeterDay := 0, 0.0, 0
	for day := 0; day < days; day++ {
		if utils.WithinPolygon(sun, planets...) {
			var positions []string
			for _, planet := range planets {
				positions = append(positions, fmt.Sprintf("r=%v, %v°", planet.R, planet.Degrees))
			}

			if perimeter := utils.Perimeter(planets...); perimeter > maxPerimeter {
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
		}

		for i := 0; i < len(planets); i++ {
			planets[i].AdvanceDay()
		}
	}

	fmt.Printf("amount of rainy periods %d\n", amountRains)
	fmt.Printf("max peak of rains at day %d\n", maxPerimeterDay)
}
