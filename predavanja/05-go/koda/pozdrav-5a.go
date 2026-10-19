// pozdrav
// dve gorutini, argument, vrstni red izpisov se od zagona do zagona spreminja

package main

import (
	"fmt"
	"sync"
	"time"
)

const printouts = 10

var wg sync.WaitGroup

func hello(s string) {
	defer wg.Done()
	for i := 0; i < printouts; i++ {
		fmt.Print(s, " ")
		time.Sleep(time.Millisecond)
	}
}

func main() {
	wg.Add(1)
	go hello("hello")
	wg.Add(1)
	go hello("world")
	wg.Wait()
	fmt.Println()
}
