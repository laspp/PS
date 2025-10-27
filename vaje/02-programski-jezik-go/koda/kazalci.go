/*
Program kazalci prikazuje uporabo kazalcev v jeziku go
*/
package main

import "fmt"

type point struct {
	x int
	y int
}

func main() {
	i, j := 42, 1337
	// Ustvarimo kazalec p, ki kaže na spremenljivko i
	p := &i
	// Dostop do vrednosti, na katero kaže kazalec
	fmt.Println(*p)
	// Spremenimo vrednost preko kazalca
	*p = j
	fmt.Println(i)

	// Pri kazalcih na strukture lahko uporabimo poenostavljeno sintakso brez *, ko spreminjamo vrednosti preko kazalca
	t := point{}
	p2 := &t
	p2.x = 15
	fmt.Println(t)
}
