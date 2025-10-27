/*
Program zakasnjeno-izvajanje prikazuje uporabo stavka defer v jeziku go
*/
package main

import "fmt"

func main() {
	// Vrednost s se bo določila takoj, izpis pa se bo zgodil na koncu
	s := "svet"
	defer fmt.Println(s)
	s = "Pozdravljen"
	fmt.Println(s)

	//Klici zakasnjenih funkcij se nalagajo na sklad in se izvedejo po principu LIFO
	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}
}
