// pozdrav
// kanal z dovolj velikim medpomnilnikom, preverimo še, kako je z medpomnilnikom velikosti 2*printouts-1

package main

import (
	"fmt"
	"strconv"
	"time"
)

const printouts = 10

var stringStream chan string

func hello(str string) {
	fmt.Println("Goroutine", str, "start")
	for i := 0; i < printouts; i++ {
		stringStream <- str + "-" + strconv.Itoa(i)
	}
	fmt.Println("Goroutine", str, "done")
}

func main() {

	stringStream = make(chan string, 2*printouts)

	go hello("hello")
	go hello("world")

	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 2*printouts; i++ {
		msg := <-stringStream
		fmt.Print(msg, " ")
	}
	fmt.Println()
}
