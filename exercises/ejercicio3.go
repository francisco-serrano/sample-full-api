package exercises

import (
	"fmt"
	"github.com/sample-full-api/models"
	"github.com/sample-full-api/utils"
	"strings"
)

func AmountOptimalConditions(days int, logSituations bool) {
	ferengi, _ := models.NewPoint(500, 45, 1, true)
	betasoide, _ := models.NewPoint(2000, 270, 3, true)
	vulcano, _ := models.NewPoint(1000, 135, 5, false)

	planets := []*models.Point{ferengi, betasoide, vulcano}

	amount := 0
	for day := 0; day < days; day++ {
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
		}

		for _, planet := range planets {
			planet.AdvanceDay()
		}
	}

	fmt.Printf("amount of optimal conditions %d\n", amount)
}
