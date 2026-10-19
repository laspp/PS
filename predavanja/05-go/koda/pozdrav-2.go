// pozdrav
// ni izpisa - glavna gorutina se prehitro zaključi

package main

import (
	"fmt"
	"time"
)

const printouts = 10

func hello() {
	for i := 0; i < printouts; i++ {
		fmt.Print("hello world ")
		time.Sleep(time.Millisecond)
	}
}

func main() {
	go hello()
	// preostala koda
	fmt.Println()
}
