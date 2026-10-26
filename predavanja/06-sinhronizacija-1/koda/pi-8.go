// računanje pi po metodi Monte Carlo
// vzporedno, se izognemo ključavnicam, prenos vrednosti s kazalcem

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type experiment struct {
	shots int				// 8B na 64-bitni ahitekturi
	hits  int				// 8B	
	value float64			// 8B
//	padding [13] int		// 24 B vs. 128 B
}

var pi experiment
var wg sync.WaitGroup

func piCalc(iterations int, result *experiment) {
	defer wg.Done()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < iterations; i++ {
		x := rnd.Float64()
		y := rnd.Float64()
		result.shots++
		if x*x+y*y < 1 {
			result.hits++
		}
	}
}

func main() {
	// preberemo argumente
	iPtr := flag.Int("i", 10000000, "# of iterations")
	gPtr := flag.Int("g", 2, "# of goroutines")
	flag.Parse()
	// dinamično alociramo rezino (slice)
	var mypi = make([]experiment, *gPtr)
	// razdelimo delo med gorutine in jih zaženemo
	timeStart := time.Now()
	wg.Add(*gPtr)
	for id := 0; id < *gPtr; id++ {
		go piCalc(*iPtr / *gPtr, &mypi[id])
	}
	// gorutine pridružimo
	wg.Wait()
	for _, mp := range mypi {
		pi.shots += mp.shots
		pi.hits += mp.hits
	}
	pi.value = 4.0 * float64(pi.hits) / float64(pi.shots)
	timeElapsed := time.Since(timeStart)
	// izpišemo
	fmt.Println("pi:", pi, "goroutines:", *gPtr, "time:", timeElapsed)
}
