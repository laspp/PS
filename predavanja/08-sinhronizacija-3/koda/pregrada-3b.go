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
var rwlock sync.RWMutex
var phase int = 0

func barrier(id int, printouts int) {
	defer wg.Done()
	var p int

	for i := 0; i < printouts; i++ {

		// operacije v zanki
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		fmt.Println("Gorutine", id, "printout", i)

		// pregrada - začetek
		// vrata 0
		rwlock.Lock()
		if g > 0 {
			rwlock.Unlock()
			p = 1
			for p == 1 {
				rwlock.Rlock()
				p = phase
				rwlock.Runlock()
			}
			rwlock.Lock()
		} else {
			phase = 0 // prehajanje čez vrata 0 se začne, ko zadnja gorutina zapusti vrata 1
		}
		g++
		rwlock.Unlock()

		// vrata 1
		rwlock.Lock()
		if g < goroutines {
			rwlock.Unlock()
			p = 0
			for p == 0 {
				rwlock.Rlock()
				p = phase
				rwlock.Runlock()
			}
			rwlock.Lock()
		} else {
			phase = 1 // prehajanje čez vrata 1 se začne, ko zadnja gorutina zapusti vrata 0
		}
		g--
		rwlock.Unlock()
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
