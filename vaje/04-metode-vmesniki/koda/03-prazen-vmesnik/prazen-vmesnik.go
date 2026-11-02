/*
Program prikazuje uporabo praznih vmesnikov v programskem jeziku Go
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

// Funkcija z izbiro tipa različne podatkovne tipe obravnava različno
func describe(x any) {
	switch v := x.(type) {
	case string:
		fmt.Println("Niz dolžine", len(v))
	case int:
		fmt.Println("Celo število", v)
	default:
		fmt.Printf("Neznan tip: %T\n", v)
	}
}

func main() {
	// Uporaba funkcije display nad različnimi podatkovnimi tipi
	s := "PS"
	display(s)
	i := 42
	display(i)
	rect := struct {
		width  float32
		height float32
	}{
		width:  5,
		height: 4,
	}
	display(rect)

	// Preverjanje tipa vrednosti, shranjene v praznem vmesniku
	var x any = s
	str, ok := x.(string)
	fmt.Println(str, ok) // PS true
	n, ok := x.(int)
	fmt.Println(n, ok) // 0 false

	// Brez druge povratne vrednosti bi naslednji stavek povzročil paniko
	// m := x.(int)

	// Uporaba funkcije describe nad različnimi podatkovnimi tipi
	describe(s)
	describe(i)
	describe(rect)
}
