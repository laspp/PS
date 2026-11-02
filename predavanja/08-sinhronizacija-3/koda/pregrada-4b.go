// Pregrada
// en kanal, vsakič naredimo nova vrata

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup
var goroutines int
var g int
var lock sync.Mutex
var door = make(chan struct{})

func barrier(id int, printouts int) {

	defer wg.Done()

	for i := 0; i < printouts; i++ {

		// operacije v zanki
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		fmt.Println("Goroutine", id, "printout", i)

		// pregrada - začetek
		lock.Lock()
		myDoor := door
		g++
		if g == goroutines {
			g = 0
			door = make(chan struct{})
			close(myDoor)
		}
		lock.Unlock()
		<-myDoor
		// pregrada - konec
	}
}

func main() {
	// preberemo argumente
	gPtr := flag.Int("g", 4, "# of goroutines")
	pPtr := flag.Int("p", 5, "# of printouts")
	flag.Parse()

	goroutines = *gPtr
	door = make(chan struct{}, goroutines)

	// zaženemo gorutine
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go barrier(i, *pPtr)
	}
	// počakamo, da vse zaključijo
	wg.Wait()
}
