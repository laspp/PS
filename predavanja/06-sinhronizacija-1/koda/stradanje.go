// Demonstracija stradanja gorutin
//  - gorutini v eni iteraciji opravita enako dela: 90 µs neodvisnega in 10 µs v kritičnem odseku
//  - greedyWorker vir zaklene za celotno iteracijo, politeWorker pa le za kritični odsek
//  - greedyWorker vir takoj po sprostitvi spet zaklene; politeWorker ga dobi šele, ko se sproži sync.Mutex 
//    logika za preprečavanje stradanja (po 1 ms)
//  - simulacija dela work in workCritical: s preverjanjem pogoja v zanki trošimo procesorske vire

package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var lock sync.Mutex

func work(d time.Duration) { t := time.Now(); for time.Since(t) < d {} }
func workCritical(d time.Duration) { t := time.Now(); for time.Since(t) < d {} }

func greedyWorker(runtime time.Duration) {
	defer wg.Done()

	count := 0
	timeStart := time.Now()
	for time.Since(timeStart) < runtime {
		lock.Lock()
		work(90 * time.Microsecond)
		workCritical(10 * time.Microsecond)
		lock.Unlock()
		count++
	}

	fmt.Println("  greedy:", count, "iterations")
}

func politeWorker(runtime time.Duration) {
	defer wg.Done()

	count := 0
	timeStart := time.Now()
	for time.Since(timeStart) < runtime {
		work(90 * time.Microsecond)
		lock.Lock()
		workCritical(10 * time.Microsecond)
		lock.Unlock()
		count++
	}

	fmt.Println("  polite:", count, "iterations")
}

func main() {
	// preberemo argumente
	runtimePtr := flag.Int("t", 1, "runtime in seconds")
	flag.Parse()
	runtime := time.Duration(*runtimePtr) * time.Second

	wg.Add(2)
	go greedyWorker(runtime)
	go politeWorker(runtime)
	wg.Wait()
}
