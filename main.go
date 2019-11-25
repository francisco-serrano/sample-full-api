package main

import "github.com/sample-full-api/exercises"

func main() {
	exercises.AmountDroughts(365*10, false)
	exercises.AmountRainyPeriods(365*10, false)
	exercises.AmountOptimalConditions(365*10, false)
}
