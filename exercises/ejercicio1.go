package exercises

import (
	"fmt"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/utils"
	"strings"
)

func AmountDroughts(days int, logSituations bool, srcPlanets ...models.Planet) {
	//ferengi, _ := models.NewPoint(500, 0, 1, true)
	//betasoide, _ := models.NewPoint(2000, 0, 3, true)
	//vulcano, _ := models.NewPoint(1000, 0, 5, false)
	//
	//solarSystem := []*models.Point{ferengi, betasoide, vulcano}

	planets := make([]models.Planet, len(srcPlanets))
	copy(planets, srcPlanets)

	amountAlignments := 0
	for day := 0; day < days; day++ {
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
		}

		for i := 0; i < len(planets); i++ {
			planets[i].AdvanceDay()
		}
	}

	fmt.Printf("amount of droughts %d\n", amountAlignments)
}
