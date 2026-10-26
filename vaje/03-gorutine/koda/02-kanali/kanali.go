/*
Program kanali demonstrira uporabo kanalov v programskem jeziku Go
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
	fmt.Println("Vsi delavci so zaključili")
	ch := make(chan int, 3)
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
		}
		// Sporočimo, da novih vrednosti ne bo več
		close(ch)
	}()
	// Zanka se zaključi, ko je kanal zaprt in prazen
	for v := range ch {
		fmt.Println(v)
	}
	// Branje iz zaprtega kanala vrne ničelno vrednost in ok == false
	v, ok := <-ch
	fmt.Println(v, ok) // 0 false
}
