// demonstracija živega objema
// dve osebi, vsaka najprej vzame svoje vilice, potem še druge
// če ne uspe dobiti drugih vilic, odloži tudi prve
// srun --reservation=fri --cpus-per-task=2 go run zivi-objem.go

package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var fork [2]sync.Mutex

func ticker() <-chan struct{} {

	signalChan := make(chan struct{})

	go func() {
		for {
			time.Sleep(1 * time.Second)
			signalChan <- struct{}{}
			signalChan <- struct{}{}
		}
	}()

	return signalChan
}

func person(signalChan <-chan struct{}, id int) {

	defer wg.Done()

	for {
		// obe gorutini čakata na signal, da spet poskusita vzeti vilice
		<-signalChan

		fork[id].Lock()
		fmt.Println("Person", id, "took fork", id)
		time.Sleep(100 * time.Millisecond)
		if fork[(id+1)%2].TryLock() {
			fmt.Println("Person", id, "took fork", id%2)
			break
		}
		fork[id].Unlock()
		fmt.Println("Person", id, "released fork", id)
		time.Sleep(100 * time.Millisecond)
	}
	fork[id].Unlock()
	fmt.Println("Person", id, "released fork", id)
	fork[(id+1)%2].Unlock()
	fmt.Println("Person", id, "released fork", id%2)
}

func main() {

	// sprožimo uro
	sigChan := ticker()

	// zaženemo gorutini
	wg.Add(2)
	go person(sigChan, 0)
	go person(sigChan, 1)

	// gorutini pridružimo
	wg.Wait()
}
