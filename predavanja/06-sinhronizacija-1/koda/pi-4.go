// računanje pi po metodi Monte Carlo
// vzporedno, uvedemo spremenljivko lock
// še vedno imamo tvegano stanje, počasno računanje

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type experiment struct {
	shots int
	hits  int
	value float64
}

var pi experiment
var wg sync.WaitGroup
var lock = 0

func piCalc(iterations int, goroutines int, id int) {
	defer wg.Done()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < iterations; i++ {
		x := rnd.Float64()
		y := rnd.Float64()
		for lock != id {
		}
		pi.shots++
		if x*x+y*y < 1 {
			pi.hits++
		}
		lock = (lock + 1) % goroutines
	}
}

func main() {
	// preberemo argumente
	iPtr := flag.Int("i", 10000000, "# of iterations")
	gPtr := flag.Int("g", 2, "# of goroutines")
	flag.Parse()
	// razdelimo delo med gorutine in jih zaženemo
	timeStart := time.Now()
	wg.Add(*gPtr)
	for id := 0; id < *gPtr; id++ {
		go piCalc(*iPtr / *gPtr, *gPtr, id)
	}
	// gorutine pridružimo
	wg.Wait()
	pi.value = 4.0 * float64(pi.hits) / float64(pi.shots)
	timeElapsed := time.Since(timeStart)
	// izpišemo
	fmt.Println("pi:", pi, "goroutines:", *gPtr, "time:", timeElapsed)
}
