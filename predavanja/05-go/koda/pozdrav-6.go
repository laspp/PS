// pozdrav
// dve dodatni gorutini, sporočili pošiljata v kanal; glavna gorutina izpisuje vse, kar se pojavi v kanalu
// če zahtevamo več podatkov, kot jih gorutine pošljejo v kanal, pride do smrtnega objema
// opazujemo, kdaj gorutine zaključijo

package main

import (
	"fmt"
	"strconv"
	"time"
)

const printouts = 10

var stringStream chan string

func hello(s string) {
	fmt.Println("Goroutine", s, "start")
	for i := 0; i < printouts; i++ {
		stringStream <- s + "-" + strconv.Itoa(i)
	}
	fmt.Println("Goroutine", s, "done")
}

func main() {

	stringStream = make(chan string)

	go hello("hello")
	go hello("world")

	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 2*printouts; i++ {
		msg := <-stringStream
		fmt.Print(msg, " ")
	}
	fmt.Println()
}
