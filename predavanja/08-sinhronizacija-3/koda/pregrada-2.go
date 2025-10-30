// Pregrada
// ključavnica
// primerna za enkratno uporabo, pri večkratni uporabi pride do smrtnega objema

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

func barrier(id int, printouts int) {
	defer wg.Done()

	for i := 0; i < printouts; i++ {

		// operacije v zanki
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		fmt.Println("Gorutine", id, "printout", i)

		// pregrada - začetek
		// * štetje gorutin
		lock.Lock()
		g++
		lock.Unlock()
		for g < goroutines {
		}
		// * ničenje števca gorutin
		if id == 0 {
			lock.Lock()
			g = 0
			lock.Unlock()
		}
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
