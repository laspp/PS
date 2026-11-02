// Pregrada
// ključavnica in pogojna spremenljivka
// wait je ovit v zanko, uporabimo zastavico za označevanje lihe/sode zanke

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
var cond = sync.NewCond(&lock)
var barrierFlag bool

func barrier(id int, printouts int) {
	defer wg.Done()

	for i := 0; i < printouts; i++ {

		// operacije v zanki
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		fmt.Println("Goroutine", id, "printout", i)

		// pregrada - začetek
		cond.L.Lock()
		myFlag := barrierFlag
		g++
		if g == goroutines {
			g = 0
			barrierFlag = !barrierFlag
			cond.Broadcast()
		}
		for myFlag == barrierFlag {
			cond.Wait()
		}
		cond.L.Unlock()
		// pregrada - konec
	}
}

func main() {
	// preberemo argumente
	gPtr := flag.Int("g", 4, "# of goroutines")
	pPtr := flag.Int("p", 5, "# of printouts")
	flag.Parse()

	goroutines = *gPtr

	// zaženemo gorutine
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go barrier(i, *pPtr)
	}
	// počakamo, da vse zaključijo
	wg.Wait()
}
