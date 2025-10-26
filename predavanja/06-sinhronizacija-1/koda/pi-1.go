// računanje pi po metodi Monte Carlo
// sekvenčno

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

type experiment struct {
	shots int
	hits  int
	value float64
}

var pi experiment

func piCalc(iterations int) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < iterations; i++ {
		x := rnd.Float64()
		y := rnd.Float64()
		pi.shots++
		if x*x+y*y < 1 {
			pi.hits++
		}
	}
}

func main() {
	// preberemo argumente
	iPtr := flag.Int("i", 10000000, "# of iterations")
	gPtr := flag.Int("g", 1, "# of goroutines")
	flag.Parse()
	// izračunamo
	timeStart := time.Now()
	piCalc(*iPtr / *gPtr)
	pi.value = 4.0 * float64(pi.hits) / float64(pi.shots)
	timeElapsed := time.Since(timeStart)
	// izpišemo
	fmt.Println("pi:", pi, "goroutines:", *gPtr, "time:", timeElapsed)
}
