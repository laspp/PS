// Demonstracija stradanja gorutin
//  - obe gorutini opravita približno enako dela
//  - gorutina greedyWorker si prilasti vir enkrat za dlje časa
//  - gorutina politeWorker si vir prilasti večkrat za krajši čas

package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var lock sync.Mutex

func politeWorker(runtime time.Duration) {
	defer wg.Done()

	count := 0
	timeStart := time.Now()
	for time.Since(timeStart) < runtime {
		lock.Lock()
		time.Sleep(1 * time.Nanosecond)
		lock.Unlock()
		lock.Lock()
		time.Sleep(1 * time.Nanosecond)
		lock.Unlock()
		lock.Lock()
		time.Sleep(1 * time.Nanosecond)
		lock.Unlock()
		count++
	}

	fmt.Println("Prijazna gorutina:", count, "iteracij.")
}

func greedyWorker(runtime time.Duration) {
	defer wg.Done()

	count := 0
	timeStart := time.Now()
	for time.Since(timeStart) < runtime {
		lock.Lock()
		time.Sleep(3 * time.Nanosecond)
		lock.Unlock()
		count++
	}

	fmt.Println("Pohlepna gorutina:", count, "iteracij.")
}

func main() {
	// preberemo argumente
	runtimePtr := flag.Int("t", 1, "runtime in seconds")
	flag.Parse()

	// zaženemo oba delavca
	wg.Add(2)
	runtime := time.Duration(*runtimePtr) * time.Second
	go politeWorker(runtime)
	go greedyWorker(runtime)
	wg.Wait()
}
