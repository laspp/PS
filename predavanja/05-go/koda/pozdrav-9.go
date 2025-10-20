// pozdrav
// dodatna gorutina zapre kanal
// glavna gorutina bere s ključno besedo range

package main

import (
	"fmt"
	"strconv"
	"time"
)

const printouts = 10

var stringStream chan string

func hello(str string) {
	defer close(stringStream)

	fmt.Println("Goroutine", str, "start")
	for i := 0; i < printouts; i++ {
		stringStream <- str + "-" + strconv.Itoa(i)
	}
	fmt.Println("Goroutine", str, "done")
}

func main() {

	stringStream = make(chan string)

	go hello("hello-world")

	time.Sleep(100 * time.Millisecond)

	for msg := range stringStream {
		fmt.Print(msg, " ")
	}
	fmt.Println()
}
