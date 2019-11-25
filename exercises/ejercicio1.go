package exercises

import (
	"fmt"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/utils"
	"strings"
)

func AmountDroughts(days int) {
	ferengi, _ := models.NewPoint(500, 0, 1, true)
	betasoide, _ := models.NewPoint(2000, 0, 3, true)
	vulcano, _ := models.NewPoint(1000, 0, 5, false)

	solarSystem := []*models.Point{ferengi, betasoide, vulcano}

	fmt.Printf("analyzing droughts in the following %d days\n", days)

	amountAlignments := 0
	for day := 0; day < days; day++ {
		if utils.AlignedWithSun(solarSystem...) {
			var positions []string
			for _, planet := range solarSystem {
				positions = append(positions, fmt.Sprintf("%v", planet.Degrees))
			}

			fmt.Printf("drought detected at day %v\t\tpositions %s\n",
				day,
				strings.Join(positions, ";"),
			)
			amountAlignments += 1
		}

		for _, planet := range solarSystem {
			planet.AdvanceDay()
		}
	}

	fmt.Printf("amount of droughts %d\n", amountAlignments)
}
