/*
Program zanke prikazuje kako uporabljamo zanko for v jeziku go
*/
package main

import "fmt"

func main() {
	// Tipičen primer zanke
	for i := 0; i <= 9; i++ {
		fmt.Println(i)
	}

	// Izpustimo inicializacijo in stavek po koncu iteracije (stavek while)
	// Podpičij nam ni treba pisati
	i := 0
	for i <= 9 {
		fmt.Println(i)
		i++
	}

	// Izpustimo še pogoj in dobimo neskončno zanko
	for {
		fmt.Println("Zanka")
		break
	}

	// Bolj elegantna oblika zanke for z uporabo range
	// Uporabno za sprehajanje po tabelah, rezinah, nizih, slovarjih
	fibonacci := [6]int{1, 1, 2, 3, 5, 8}
	for i, v := range fibonacci {
		fmt.Println(i, v)
	}
	fmt.Println()
	// Indeks ali vrednost lahko izpustimo
	for _, v := range fibonacci {
		fmt.Println(v)
	}
	fmt.Println()
	for i := range fibonacci {
		fmt.Println(i)
	}

	// Pri nizih obstaja razlika
	// Če se po nizu sprehodimo z range nam vrača rune (znaki unicode)
	// Če se po nizu sprehodimo z indeksom nam vrača posamezne bajte
	s := "Čuri Muri, kljukec Juri"
	for _, c := range s {
		fmt.Printf("%c", c)
	}
	fmt.Println()
	for i := 0; i < len(s); i++ {
		fmt.Printf("%c", s[i])
	}
	fmt.Println()
}
