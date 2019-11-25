package main

import "github.com/sample-full-api/exercises"

func main() {
	exercises.AmountDroughts(365 * 10)
	exercises.AmountRainyPeriods(365 * 10)

	//ferengi, _ := models.NewPoint(500, 45, 1, true)
	//betasoide, _ := models.NewPoint(2000, 270, 3, true)
	//vulcano, _ := models.NewPoint(1000, 135, 5, false)
	//
	//fmt.Println("perimeter", utils.Perimeter(ferengi, betasoide, vulcano))

	//a, _ := models.NewPoint(10, 90)
	//b, _ := models.NewPoint(20, 270)
	//c, _ := models.NewPoint(18, 0)
	//d, _ := models.NewPoint(15, 90)
	//e, _ := models.NewPoint(30, 270)
	//
	//points := []*models.Point{a, b, c, d, e}
	//
	//for _, p := range points {
	//	fmt.Printf("%+v\n", p)
	//}
	//
	//f, _ := models.NewPoint(10, 180)
	//
	//fmt.Printf("sample point %+v\n", f)
	//
	//fmt.Println("aligned with sun", AlignedWithSun(points...))
	//fmt.Printf("determinant\na\t\t%+v\nb\t\t%+v\ntarget\t%+v\ndet\t\t%f\n", a, b, c, Determinant(a, b, c))
	//
	//g, _ := models.NewPoint(7.07, 45)  // x=5 y=5
	//h, _ := models.NewPoint(7.07, 135) // x=-5 y=5
	//i, _ := models.NewPoint(5, 270)    // x=0 y=-5
	//j, _ := models.NewPoint(7.07, 315) // x=5 y=-5
	//
	//target, _ := models.NewPoint(7.07, 225)
	//
	//fmt.Printf("%+v\n", g)
	//fmt.Printf("%+v\n", h)
	//fmt.Printf("%+v\n", i)
	//fmt.Printf("%+v\n", j)
	//fmt.Printf("%+v\n", target)
	//
	//fmt.Println(Determinant(g, h, target))
	//fmt.Println(Determinant(h, i, target))
	//fmt.Println(Determinant(i, j, target))
	//fmt.Println(Determinant(j, g, target))
	//
	//polygon := []*models.Point{g, h, i, j}
	//
	//fmt.Println(WithinPolygon(target, polygon...))
	//
	//fmt.Println(AlignedWithoutSun(a, b, d, e))
}
