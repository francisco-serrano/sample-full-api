package models

import (
	"errors"
	"math"
)

type Point struct {
	r       float64
	Degrees float64
	radians float64
	X       float64
	Y       float64
}

func NewPoint(r, degrees float64) (*Point, error) {
	if r < 0.0 || degrees < 0.0 || degrees >= 360.0 {
		return nil, errors.New("invalid input data")
	}

	radians := degrees * (math.Pi / 180.0)

	return &Point{
		r:       r,
		Degrees: degrees,
		radians: radians,
		X:       math.Round(r*math.Cos(radians)*100) / 100,
		Y:       math.Round(r*math.Sin(radians)*100) / 100,
	}, nil
}

