/*
Program select prikazuje uporabo stavka select pri delu s kanali v programskem jeziku go
*/
package main

import (
	"fmt"
	"time"
)

// Funkcija, ki čaka na pritisk tipke "Enter"
func readKey(input chan bool) {
	fmt.Scanln()
	input <- true
}

// Delavec
func worker(id int, done chan bool) {
	fmt.Print("Delavec ", id)
	time.Sleep(2 * time.Second)
	fmt.Print("Končal")
	done <- true
}

func main() {
	input := make(chan bool)
	done := make(chan bool)
	w := 0
	// Zaženemo gorutino, ki čaka na pritisk tipke
	go readKey(input)
	// Zaženemo prvega delavca
	go worker(w, done)
	// Anonimna funkcija z neskončno zanko
	func() {
		for {
			select {
			// Pritisk tipke: končamo
			case <-input:
				return
			// Delavec zaključil, zaženimo novega
			case <-done:
				fmt.Println()
				w = w + 1
				go worker(w, done)
			// Če se nič ne zgodi izvedemo privzeto akcijo
			default:
				time.Sleep(200 * time.Millisecond)
				fmt.Print(".")
			}
		}
	}()
	// Počakamo na zadnjega delavca
	<-done
}
