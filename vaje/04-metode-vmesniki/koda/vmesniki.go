/*
Program prikazuje uporabo vmesnikov v programskem jeziku Go
*/
package main

import (
	"fmt"
	"math"
)

// Definiramo nov vmesnik za izračun površine lika
type areaCalculator interface {
	area() float32
}

// Definirmao strukturo za krog
type circle struct {
	radius float32
}

// Definiramo strukturo za pravokotnik
type rect struct {
	width  float32
	heigth float32
}

// Metoda ta izračun površine kroga
func (c circle) area() float32 {
	return c.radius * c.radius * math.Pi

}

// Metoda ta izračun površine pravokotnika
func (r rect) area() float32 {
	return r.width * r.heigth
}

func main() {

	// Ustvarimo nekaj likov
	circle1, circle2, circle3 := circle{10}, circle{3}, circle{8}
	rect1, rect2, rect3 := rect{2, 2}, rect{3, 8}, rect{10, 10}

	// Ustvarimo rezino vmesnikov in vanjo damo vse like, ki implementirajo vmesnik
	shapes := []areaCalculator{circle1, circle2, circle3, rect1, rect2, rect3}

	// Izračunajmo skupno površino vseh likov v rezini
	totalArea := float32(0.0)
	for _, v := range shapes {
		totalArea = totalArea + v.area()
	}
	fmt.Printf("Skupna površina likov: %f", totalArea)
}
