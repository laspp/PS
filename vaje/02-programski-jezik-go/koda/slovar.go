/*
Program slovar prikazuje kako definiramo, inicializiramo in dostopamo do podatkovne strukture slovar (map) v programskem jeziku go
*/
package main

import "fmt"

func main() {
	// Definiramo in inicializiramo nov slovar
	var fruitSupply map[string]int
	// Dokler slovarja ne inicializiramo z make, do njega ne moremo dostopati
	fruitSupply = make(map[string]int)

	// Dodamo nove ključe v slovar
	fruitSupply["jabolka"] = 5
	fruitSupply["hruske"] = 3
	fmt.Println(fruitSupply)
	// Za ključe, ki še niso v slovarju dobimo ničelno vrednost
	fmt.Println(fruitSupply["jabolka"], fruitSupply["slive"])

	// Izbrišemo ključ
	delete(fruitSupply, "jabolka")

	//Preverimo, če je ključ v slovarju
	e, ok := fruitSupply["hruske"]
	fmt.Println(e, ok)

	// Pobrišemo vse ključe
	clear(fruitSupply)
	fmt.Println(fruitSupply)

	// Ustvarimo nov slovar z zalogo sadja z neposredno inicializacijo ključev
	fruitSupply = map[string]int{"jabolka": 1, "hruske": 2, "slive": 10}
	fmt.Println(fruitSupply)
}
