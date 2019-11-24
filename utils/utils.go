package utils

import (
	"github.com/sample-full-api/models"
	"sort"
)

func AlignedWithoutSun(points ...*models.Point) bool {
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

func AlignedWithSun(points ...*models.Point) bool {
	if len(points) < 3 {
		return alignedWithSun(points[0], points[1])
	}

	for i := 0; i < len(points)-1; i++ {
		if !alignedWithSun(points[i], points[i+1]) {
			return false
		}
	}

	return true
}

func alignedWithSun(a, b *models.Point) bool {
	return a.Degrees == b.Degrees || (a.Degrees+180.0 == b.Degrees) || (a.Degrees == b.Degrees+180.0)
}

// polygon must be convex
func WithinPolygon(target *models.Point, polygon ...*models.Point) bool {
	if len(polygon) < 3 {
		return false
	}

	sortedPolygon := make([]*models.Point, len(polygon))
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

func Determinant(a, b, target *models.Point) float64 {
	return (b.X-a.X)*(target.Y-a.Y) - (b.Y-a.Y)*(target.X-a.X)
}
