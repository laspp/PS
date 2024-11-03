// Problem pisateljev in bralcev
// uporabimo ključavnice, bralce štejemo

package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var activeReaders int = 0
var lockReaders sync.Mutex
var lockBook sync.Mutex

func writer(id int, cycles int) {
	defer wg.Done()

	for i := 0; i < cycles; i++ {
		lockBook.Lock()
		fmt.Println("Writer", id, "start", i)
		time.Sleep(time.Duration(id) * time.Millisecond)
		fmt.Println("Writer", id, "finish", i)
		lockBook.Unlock()
		time.Sleep(time.Duration(id) * time.Millisecond)
	}
}

func reader(id int) {

	for {
		lockReaders.Lock()
		activeReaders++
		if activeReaders == 1 {
			lockBook.Lock()
		}
		lockReaders.Unlock()

		fmt.Println("Reader", id, "start")
		time.Sleep(time.Duration(id) * time.Millisecond)
		fmt.Println("Reader", id, "finish")

		lockReaders.Lock()
		activeReaders--
		if activeReaders == 0 {
			lockBook.Unlock()
		}
		lockReaders.Unlock()

		time.Sleep(time.Duration(id) * time.Millisecond)
	}
}

func main() {
	// preberemo argumente
	writersPtr := flag.Int("w", 2, "# of writers")
	readersPtr := flag.Int("r", 4, "# of readers")
	cyclesPtr := flag.Int("c", 10, "# of cycles")
	flag.Parse()

	// zaženemo pisatelje
	wg.Add(*writersPtr)
	for i := 1; i <= *writersPtr; i++ {
		go writer(i, *cyclesPtr)
	}
	// zaženemo bralce
	for i := 1; i <= *readersPtr; i++ {
		go reader(i)
	}
	// počakamo, da pisatelji zaključijo
	wg.Wait()
}
