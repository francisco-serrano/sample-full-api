package utils

import (
	"github.com/sample-full-api/models"
	"math"
	"sort"
)

func AlignedWithoutSun(points ...models.Planet) bool {
	if len(points) < 3 {
		return true
	}

	for i := 0; i < len(points)-2; i++ {
		if Determinant(points[i], points[i+1], points[i+2]) != 0 {
			return false
		}
	}

	return true
}

func AlignedWithSun(planets ...models.Planet) bool {
	if len(planets) < 3 {
		return alignedWithSun(planets[0], planets[1])
	}

	for i := 0; i < len(planets)-1; i++ {
		if !alignedWithSun(planets[i], planets[i+1]) {
			return false
		}
	}

	return true
}

func alignedWithSun(a, b models.Planet) bool {
	return a.Degrees == b.Degrees || (a.Degrees+180.0 == b.Degrees) || (a.Degrees == b.Degrees+180.0)
}

// polygon must be convex
func WithinPolygon(target models.Planet, polygon ...models.Planet) bool {
	if len(polygon) < 3 {
		return false
	}

	sortedPolygon := make([]models.Planet, len(polygon))
	copy(sortedPolygon, polygon)

	sort.Slice(sortedPolygon, func(i, j int) bool {
		return polygon[i].Degrees < polygon[j].Degrees
	})

	for i := 0; i < len(sortedPolygon)-1; i++ {
		if Determinant(sortedPolygon[i], sortedPolygon[i+1], target) < 0 {
			return false
		}
	}

	if Determinant(sortedPolygon[len(sortedPolygon)-1], sortedPolygon[0], target) < 0 {
		return false
	}

	return true
}

func Determinant(a, b, target models.Planet) float64 {
	return (b.X-a.X)*(target.Y-a.Y) - (b.Y-a.Y)*(target.X-a.X)
}

func Perimeter(polygon ...models.Planet) float64 {
	if len(polygon) < 2 {
		return 0
	}

	if len(polygon) == 2 {
		return distance(polygon[0], polygon[1])
	}

	perimeter := 0.0
	for i := 0; i < len(polygon)-1; i++ {
		perimeter += distance(polygon[i], polygon[i+1])
	}

	perimeter += distance(polygon[len(polygon)-1], polygon[0])

	return perimeter
}

func distance(a, b models.Planet) float64 {
	return math.Sqrt(math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2))
}
