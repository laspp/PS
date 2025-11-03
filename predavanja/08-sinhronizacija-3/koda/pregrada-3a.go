// Pregrada
// ključavnica, princip dvojih vrat:
// 		faza (phase) = 0: prehajanje čez vrata 0
//		faza 		 = 1: prehajanje čez vrata 1
//		g				: število gorutin med vrati
// tveganega stanja ni več

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
var g int = 0
var lock sync.Mutex
var phase int = 1 // vrata 0 so zaprta, odpre jih prva gorutina

func barrier(id int, printouts int) {
	defer wg.Done()

	for i := 0; i < printouts; i++ {

		// operacije v zanki
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		fmt.Println("Gorutine", id, "printout", i)

		// pregrada - začetek
		// vrata 0
		lock.Lock()
		if g > 0 {
			lock.Unlock()
			for phase == 1 {
			}
			lock.Lock()
		} else {
			phase = 0 // prehajanje čez vrata 0 se začne, ko zadnja gorutina zapusti vrata 1
		}
		g++
		lock.Unlock()

		// vrata 1
		lock.Lock()
		if g < goroutines {
			lock.Unlock()
			for phase == 0 {
			}
			lock.Lock()
		} else {
			phase = 1 // prehajanje čez vrata 1 se začne, ko zadnja gorutina zapusti vrata 0
		}
		g--
		lock.Unlock()
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
