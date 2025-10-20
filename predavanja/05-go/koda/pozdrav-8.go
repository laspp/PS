// pozdrav
// dodatna gorutina zapre kanal
// glavna gorutina bere iz kanala, gledamo zastavico OK

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

	for i := 0; i < printouts+1; i++ {
		msg, ok := <-stringStream
		fmt.Print(msg, "(", ok, ") ")
	}
	fmt.Println()
}
