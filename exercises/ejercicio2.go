package exercises

import (
	"fmt"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/utils"
	"strings"
)

func AmountRainyPeriods(days int) {
	ferengi, _ := models.NewPoint(500, 45, 1, true)
	betasoide, _ := models.NewPoint(2000, 270, 3, true)
	vulcano, _ := models.NewPoint(1000, 135, 5, false)

	planets := []*models.Point{ferengi, betasoide, vulcano}

	sun, _ := models.NewPoint(0, 0, 0, false)

	fmt.Println("within polygon", utils.WithinPolygon(sun, planets...))

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

			fmt.Printf("rainy period at day %v\t\t%s\n",
				day,
				strings.Join(positions, "\t"),
			)
			amountRains += 1
		}

		for _, planet := range planets {
			planet.AdvanceDay()
		}
	}

	fmt.Printf("amount of rainy periods %d\n", amountRains)
	fmt.Printf("max peak of rains at day %d\n", maxPerimeterDay)
}
