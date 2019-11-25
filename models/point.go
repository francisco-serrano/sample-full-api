package models

import (
	"errors"
	"math"
)

type Point struct {
	R         float64
	Degrees   float64
	radians   float64
	X         float64
	Y         float64
	Speed     float64
	clockwise bool
}

func NewPoint(r, degrees, speed float64, clockwise bool) (*Point, error) {
	if r < 0.0 || degrees < 0.0 || degrees >= 360.0 {
		return nil, errors.New("invalid input data")
	}

	radians := degrees * (math.Pi / 180.0)

	return &Point{
		R:         r,
		Degrees:   degrees,
		radians:   radians,
		X:         math.Round(r*math.Cos(radians)*100) / 100,
		Y:         math.Round(r*math.Sin(radians)*100) / 100,
		Speed:     speed,
		clockwise: clockwise,
	}, nil
}

func (p *Point) AdvanceDay() {
	if !p.clockwise {
		newPosition := p.Degrees + p.Speed

		for newPosition >= 360 {
			newPosition -= 360
		}

		p.Degrees = newPosition
		p.radians = newPosition * (math.Pi / 180.0)
		p.X = math.Round(p.R*math.Cos(p.radians)*100) / 100
		p.X = math.Round(p.R*math.Sin(p.radians)*100) / 100

		return
	}

	newPosition := p.Degrees - p.Speed

	for newPosition < 0 {
		newPosition += 360
	}

	p.Degrees = newPosition
	p.radians = newPosition * (math.Pi / 180.0)
	p.X = math.Round(p.R*math.Cos(p.radians)*100) / 100
	p.X = math.Round(p.R*math.Sin(p.radians)*100) / 100
}
