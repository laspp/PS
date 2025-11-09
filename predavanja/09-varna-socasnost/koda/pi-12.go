// računanje pi po metodi Monte Carlo
// vzporedno, rešitev s kanali
// uporabljamo samo eno strukturo rnd tipa *rand.Rand, dostop do nje zaščitimo s ključavnico
// seme je fiksna vrednost, rezultati niso ponovljivi, saj zaradi dveh kritičnih odsekov
// določamo različne točke (x, y)

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
var rnd *rand.Rand
var lock sync.Mutex

func piCalc(iterations int, result chan experiment) {
	var mypi experiment
	for i := 0; i < iterations; i++ {
		lock.Lock()
		x := rnd.Float64()
		lock.Unlock()
		lock.Lock()
		y := rnd.Float64()
		lock.Unlock()
		mypi.shots++
		if x*x+y*y < 1 {
			mypi.hits++
		}
	}
	result <- mypi
}

func main() {
	// preberemo argumente
	iPtr := flag.Int("i", 10000000, "# of iterations")
	gPtr := flag.Int("g", 2, "# of goroutines")
	sPtr := flag.Int64("s", 0, "random seed")
	flag.Parse()
	// odpremo kanal
	piStream := make(chan experiment)
	rnd = rand.New(rand.NewSource(*sPtr))
	// razdelimo delo med gorutine in jih zaženemo
	timeStart := time.Now()
	for id := 0; id < *gPtr; id++ {
		go piCalc(*iPtr / *gPtr, piStream)
	}
	// preberemo vse vrednosti iz kanala
	for i := 0; i < *gPtr; i++ {
		mypi := <-piStream
		pi.shots += mypi.shots
		pi.hits += mypi.hits
	}
	pi.value = 4.0 * float64(pi.hits) / float64(pi.shots)
	timeElapsed := time.Since(timeStart)
	// izpišemo
	fmt.Println("pi:", pi, "goroutines:", *gPtr, "time:", timeElapsed)
}
