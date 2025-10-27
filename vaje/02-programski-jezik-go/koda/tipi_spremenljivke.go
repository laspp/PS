/*
Program tipi_spremenljivke prikazuje načine definiranja različnih osnovnih tipov spremenljivk v go.
*/
package main

import "fmt"

var b bool = true
var x int = 42
var y float64 = 3.141592653
var c complex128 = 42 + 3.14159i
var ch rune = 'Č'
var s string = "PS"

func main() {
	// v isti vrstici lahko definiramo več spremenljivk istega tipa
	var i, j int = 1, 2

	// operator := omogoča hkratno definiranje spremenljivke brez ključne besede var
	// podatkovni tip se določi samodejno glede na vrednost, ki jo prirejamo
	// dovoljeno samo znotraj funkcij
	k := 3
	banana, jabolko, brokoli := true, false, "Ne!"

	//Println samodejno uporabi ustrezno formatiranje glede na podatkovni tip
	fmt.Println(b, x, y, c, ch, s)
	fmt.Printf("%c\n", ch)

	fmt.Println(i, j, k, c, banana, jabolko, brokoli)
	// %T izpiše podatkovni tip posamezne spremenljivke
	fmt.Printf("%T %T %T %T %T %T %T %T\n", ch, i, j, k, c, banana, jabolko, brokoli)

	// Med podatkovnimi tipi lahko tudi pretvarjamo
	// Za razliko od programskega jezika C, je v go vedno potrebna eksplicitna konverzija pri priredbah med različnimi podatkovnimi tipi
	fx := float64(x)
	fmt.Printf("%T, %.1f\n", fx, fx)
}
