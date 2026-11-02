// Pregrada
// dva kanala, princip dvojih vrat
// 		najprej počakamo, da vse gorutine pridejo čez vrata 0; zadnja odpre vrata 1
//      nato počakamo, da vse gorutine pridejo čez vrata 1; zadnja odpre vrata 0
// 		kanala imata kapaciteto 1, da zadnja gorutina lahko sprosti sama sebe
//		kanala sta tipa struct{}, s čimer poudarimo, da ju uporabljamo za sinhronizacijo

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
var door1 = make(chan struct{}, 1)
var door0 = make(chan struct{}, 1)

func barrier(id int, printouts int) {

	defer wg.Done()

	for i := 0; i < printouts; i++ {

		// operacije v zanki
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		fmt.Println("Goroutine", id, "printout", i)

		// pregrada - začetek
		// vrata 0 odprta
		lock.Lock()
		g++
		if g == goroutines {
			for j := 0; j < goroutines; j++ {
				door1 <- struct{}{} // zadnja gorutina odpre vrata 1
			}
		}
		lock.Unlock()
		<-door1 // gorutine čakajo pred zaprtimi vrati

		// vrata 1 odprta
		lock.Lock()
		g--
		if g == 0 {
			for j := 0; j < goroutines; j++ {
				door0 <- struct{}{} // zadnja gorutina odpre vrata 0
			}
		}
		lock.Unlock()
		<-door0 // gorutine čakajo pred zaprtimi vrati
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
