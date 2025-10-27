/*
Program funckije demonstrira uporabo funkcij v go.
Prikazani so različni načini podajanja argumentov, vračanja vrednosti, anonimne funkcije in zaprtja (closures)
*/
package main

import "fmt"

func sestej(a int, b int) int {
	return a + b
}

// Če ima več argumentov isti podatkovni tip, ga lahko pišemo samo enkrat, na koncu seznama argumentov
// Funkcija lahko vrača več vrednosti, podprto je tudi poimenovanje vračanih vrednosti
func deli(a, b int) (q, r int) {
	q = a / b
	r = a % b
	return
}

// Zaprtja funkcij (closures) omogočajo, da lahko znotraj funkcije naslavljamo spremenljivko, ki je definirana izven funkcije.
// Vrednost take spremenljive se ohranja tudi potem ko se iz funkcije vrnemo.
// Funkcija povecaj vrača anonimno funkcijo (zaprtje), ki spreminja spremenljivko i
func povecaj() func() {
	var i int
	return func() {
		i++
		fmt.Println(i)
	}
}

func main() {
	fmt.Println(sestej(42, 7))
	fmt.Println(deli(42, 7))
	// Go podpira anonimne funkcije, torej take, ki nimajo imena
	// Če anonimni funkciji sledi (), se le-ta takoj izvede
	result := func() int {

		fmt.Println("Anonimna funkcija")
		return 42
	}()
	fmt.Println(result)
	// Če pa ne, se spremenljivki priredi funkcija, ki jo lahko kasneje pokličemo
	f := func() {

		fmt.Println("Prirejena anonimna funkcija")
	}
	f()
	// Klic funckije, ki vrača zaprtje, vrednost i se skozi klice ohranja.
	p := povecaj()
	p()
	p()
}
