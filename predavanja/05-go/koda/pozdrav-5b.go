// pozdrav
// dve gorutini, argument, vrstni red izpisov se od zagona do zagona spreminja
// glavna gorutina tudi dela
// trikrat kličemo wg.Done(), zato rabimo wg.Add(1) tudi za glavno gorutino

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
	wg.Add(1)
	hello("!")
	wg.Wait()
	fmt.Println()
}
