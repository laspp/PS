// Varna sočasnost
//		z bralno-pisalno ključavnico preprečimo tvegano stanje

package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var lock sync.RWMutex

func writeToMap(id int, steps int, dict map[int]int) {
	defer wg.Done()
	for i := 0; i < steps; i++ {
		lock.Lock()
		dict[id] = i
		lock.Unlock()
	}
}

func readFromMap(id int, steps int, dict map[int]int) {
	defer wg.Done()
	for i := 0; i < steps; i++ {
		lock.RLock()
		_ = dict[id]
		lock.RUnlock()
	}
}

func main() {
	gwPtr := flag.Int("gw", 1, "# of writing goroutines")
	grPtr := flag.Int("gr", 1, "# of reading goroutines")
	sPtr := flag.Int("s", 100, "# of read or write steps")
	flag.Parse()

	dict := make(map[int]int)
	records := max(*gwPtr, *grPtr)
	for i := 0; i < records; i++ {
		dict[i] = 0
	}

	timeStart := time.Now()
	for i := 0; i < *gwPtr; i++ {
		wg.Add(1)
		go writeToMap(i, *sPtr, dict)
	}
	for i := 0; i < *grPtr; i++ {
		wg.Add(1)
		go readFromMap(i, *sPtr, dict)
	}
	wg.Wait()

	fmt.Println("dict:", dict, "time:", time.Since(timeStart))
}
