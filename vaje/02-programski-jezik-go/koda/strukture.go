/*
Program strukture prikazuje kako definiramo, inicializiramo in dostopamo do struktur v programskem jeziku go
*/
package main

import "fmt"

// Definicija strukture
type circle struct {
	x      int
	y      int
	radius int
	color  string
}

func main() {
	// Strukturo lahko inicializiramo na različne načine
	var smallCircle circle
	smallCircle.x = 0
	smallCircle.y = 0
	smallCircle.radius = 5
	smallCircle.color = "zelena"
	fmt.Println(smallCircle.x, smallCircle.y, smallCircle.radius, smallCircle.color)

	// Posamezna polja strukture lahko inicializiramo neposredno s pomočjo notacije {}
	bigCircle := circle{100, 100, 50, "rdeča"}
	fmt.Println(bigCircle)

	// Pri inicializaciji lahko navedemo imena polj, ki jim nastavljamo vrednost
	// ostala polja dobijo privzeto ničelno vrednost za dan podatkovni tip
	var mediumCircle = circle{radius: 15, color: "modra"}
	fmt.Printf("%T", mediumCircle)
}
