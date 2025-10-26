// računanje pi po metodi Monte Carlo
// vzporedno, uporabimo atomarno seštevanje

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type experimentAtomic struct {
	shots atomic.Int64
	hits  atomic.Int64
	value float64
}

var pi experimentAtomic
var wg sync.WaitGroup

func piCalc(iterations int) {
	defer wg.Done()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < iterations; i++ {
		x := rnd.Float64()
		y := rnd.Float64()
		pi.shots.Add(1)
		if x*x+y*y < 1 {
			pi.hits.Add(1)
		}
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
		go piCalc(*iPtr / *gPtr)
	}
	// gorutine pridružimo
	wg.Wait()
	pi.value = 4.0 * float64(pi.hits.Load()) / float64(pi.shots.Load())
	timeElapsed := time.Since(timeStart)
	// izpišemo
	fmt.Println("pi:", pi, "goroutines:", *gPtr, "time:", timeElapsed)
}
