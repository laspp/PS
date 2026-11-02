/*
Program prikazuje uporabo vmesnikov v programskem jeziku Go
*/
package main

import (
	"errors"
	"fmt"
	"math"
)

// Definiramo nov vmesnik za izračun površine lika
// Če površine ni mogoče izračunati, metoda vrne napako
type areaCalculator interface {
	area() (float32, error)
}

// Definiramo strukturo za krog
type circle struct {
	radius float32
}

// Definiramo strukturo za pravokotnik
type rect struct {
	width  float32
	height float32
}

// Metoda za izračun površine kroga
func (c circle) area() (float32, error) {
	if c.radius < 0 {
		return 0, errors.New("polmer kroga ne sme biti negativen")
	}
	return c.radius * c.radius * math.Pi, nil
}

// Metoda za izračun površine pravokotnika
func (r rect) area() (float32, error) {
	if r.width < 0 || r.height < 0 {
		return 0, errors.New("stranici pravokotnika ne smeta biti negativni")
	}
	return r.width * r.height, nil
}

func main() {

	// Ustvarimo nekaj likov
	circle1, circle2, circle3 := circle{10}, circle{3}, circle{8}
	rect1, rect2, rect3 := rect{2, 2}, rect{3, 8}, rect{10, 10}
	// Ustvarimo še lik z neveljavnimi podatki
	invalid := circle{-1}

	// Ustvarimo rezino vmesnikov in vanjo damo vse like, ki implementirajo vmesnik
	shapes := []areaCalculator{circle1, circle2, circle3, rect1, rect2, rect3, invalid}

	// Izračunajmo skupno površino vseh likov v rezini
	totalArea := float32(0.0)
	for _, v := range shapes {
		a, err := v.area()
		// Like, katerih površine ni mogoče izračunati, preskočimo
		if err != nil {
			fmt.Println("Napaka:", err)
			continue
		}
		totalArea = totalArea + a
	}
	fmt.Printf("Skupna površina likov: %f\n", totalArea)
}
