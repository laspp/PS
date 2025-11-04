/*
Program kanali demonstrira uporabo kanalov v programskem jeziku go
*/
package main

import (
	"fmt"
	"time"
)

func worker(id int, done chan bool) {
	fmt.Println(id, "Delam ...")
	time.Sleep(time.Second)
	fmt.Println(id, "Zaključil")

	done <- true
}

func main() {
	// Ustvarimo kanal s kapaciteto 3
	workers := 3
	done := make(chan bool, workers)
	// Zaženemo delavce
	for w := 0; w < workers; w++ {
		go worker(w, done)
	}
	// Počakamo, da delavci zaključijo
	for w := 0; w < workers; w++ {
		<-done
	}
}
