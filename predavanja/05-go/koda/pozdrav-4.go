// pozdrav
// sledimo modelu razcepi-združi

package main

import (
	"fmt"
	"time"
	"sync"
)

const printouts = 10

var wg sync.WaitGroup

func hello() {
	defer wg.Done()
	for i := 0; i < printouts; i++ {
		fmt.Print("hello world ")
		time.Sleep(time.Millisecond)		
	}
}

func main() {
	wg.Add(1)
	go hello()

	wg.Wait()
	fmt.Println()
}