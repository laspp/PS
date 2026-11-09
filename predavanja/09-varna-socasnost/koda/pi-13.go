// računanje pi po metodi Monte Carlo
// vzporedno, rešitev s kanali
// vsaka gorutina ima lokalen generator psevdonaključnih števil
// semena dodatnih gorutin izpeljemo iz generatorja nključnih števil v glavni gorutini
// pri fiksnem številu gorutin so rezultati ponovljivi

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

func piCalc(iterations int, seed int64, result chan experiment) {
	rnd := rand.New(rand.NewSource(seed))
	var mypi experiment
	for i := 0; i < iterations; i++ {
		x := rnd.Float64()
		y := rnd.Float64()
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
	// generator naključnih števil
	rndMaster := rand.New(rand.NewSource(*sPtr))
	// razdelimo delo med gorutine in jih zaženemo
	timeStart := time.Now()
	for id := 0; id < *gPtr; id++ {
		go piCalc(*iPtr / *gPtr, rndMaster.Int63(), piStream)
	}
	// preberemo vse vrednosti iz kanala (vrstni red zaključevanja se spreminja)
	for i := 0; i < *gPtr; i++ {
		mypi := <-piStream
		fmt.Println("mypi:", mypi)
		pi.shots += mypi.shots
		pi.hits += mypi.hits
	}
	pi.value = 4.0 * float64(pi.hits) / float64(pi.shots)
	timeElapsed := time.Since(timeStart)
	// izpišemo
	fmt.Println("pi:", pi, "goroutines:", *gPtr, "time:", timeElapsed)
	// zapremo kanal
	close(piStream)
}
