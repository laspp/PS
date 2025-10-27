/*
Program rezine prikazuje kako definiramo, inicializiramo in dostopamo do rezin v programskem jeziku go
*/
package main

import "fmt"

func main() {

	// fibonaci je tabela (array)
	fibonacci := [6]int{1, 1, 2, 3, 5, 8}
	// s1 je rezina (slice), ki služi kot referenca na del tabele fibonacii
	// z notacijo [:] povemo, kateri del tabele sestavlja rezino
	// nismo ustvarili kopije podatkov ampak samo referenco na elemente tabele
	var s1 []int = fibonacci[0:3]
	// Do elementov v rezini dostopamo na enak način kot do elementov tabele
	fmt.Printf("%T, %d, %d, %d\n", s1, s1[0], s1[1], s1[2])

	// Zgornjo ali spoodnjo ali obe meji lahko izpustimo, če ustvarjamo rezino, ki vsebuje vse elemente polja
	var s2 []int = fibonacci[:]
	s2[0] = 0
	fmt.Println(s2)
	fmt.Println(fibonacci)

	// Rezino lahko ustvarimo tudi neposredno
	// hkrati se ustvari tabela s podatki in rezina, ki kaže na tabelo
	letters := []string{"a", "b", "c", "d"}
	fmt.Println(letters)

	// Vsaka rezina ima definirano dolžino in kapaciteto
	fmt.Printf("Dolzina s1:%d, Kapaciteta s1:%d\nDolzina s2:%d, Kapaciteta s2:%d\n", len(s1), cap(s1), len(s2), cap(s2))

	// Dolžino rezine lahko povečamo, če je na voljo dovolj kapacitete
	// Pri tem ne pride do odvečnega kopiranja podatkov
	s1 = s1[:4]
	fmt.Printf("Dolzina s1:%d | %d, %d, %d %d\n", len(s1), s1[0], s1[1], s1[2], s1[3])

	// Rezine lahko dinamično ustvarjamo s pomočjo funkcije make
	// Rezina s3 bo inicializirana z ničelnimi vrednosti in bo imela dolžino ter kapaciteto enako 5
	s3 := make([]int, 5)
	fmt.Println(s3)

	// Dodatno lahko specificiramo tudi kapaciteto rezine
	s4 := make([]int, 5, 10)
	fmt.Printf("Dolzina s4:%d, Kapaciteta s4:%d\n", len(s4), cap(s4))

	// Poznamo tudi večdimenzionalne rezine
	magic := [][]int{
		[]int{2, 7, 6},
		[]int{9, 5, 1},
		[]int{4, 3, 8},
	}
	fmt.Println(magic)

	// Rezini lahko dodajamo elemete s funkcijo append
	// Dolžina in kapaciteta rezine se povečata
	fmt.Printf("Dolzina letters:%d, Kapaciteta letters:%d | %v\n", len(letters), cap(letters), letters)
	letters = append(letters, "e", "f")
	fmt.Printf("Dolzina letters:%d, Kapaciteta letters:%d | %v\n", len(letters), cap(letters), letters)

	// Rezine lahko kopiramo
	copy(s3, s1)
	fmt.Printf("Dolzina s3:%d, Kapaciteta s3:%d |%v -> %v\n", len(s3), cap(s3), s1, s3)
}
