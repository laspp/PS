/*
Program prikazuje uporabo praznih vmesnikov v programskem jeziku go
*/
package main

import (
	"fmt"
)

// Funkcija kot argument prejme prazen vmesnik
// Namesto interface{} lahko pišemo tudi any
func display(x interface{}) {
	fmt.Printf("Tip: %T, Vrednost: %v\n", x, x)
}

func main() {
	// Uporaba funkcije display nad različnimi podatkovnimi tipi
	s := "PS"
	display(s)
	i := 42
	display(i)
	rect := struct {
		width  float32
		heigth float32
	}{
		width:  5,
		heigth: 4,
	}
	display(rect)
}
