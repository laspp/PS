/*
Program gorutine demonstrira uporabo gorutin v jeziku go
*/
package main

import (
	"fmt"
	"time"
)

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}

func main() {
	// Zaženemo funkcijo f znotraj glavne gorutine
	f("Glavna nit")

	// Ustvarimo novo gorutino v kateri se izvaja funkcija f
	go f("Prva gorutina")

	// Ustvarimo novo gorutino, ki izvede anonimno funkcijo
	go func(msg string) {
		fmt.Println(msg)
	}("Druga gorutina")

	// Počakamo, da se vse gorutine zaključijo
	time.Sleep(time.Second)
	fmt.Println("KONEC")
}
