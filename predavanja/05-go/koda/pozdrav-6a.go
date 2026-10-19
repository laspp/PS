// pozdrav
// dve dodatni gorutini, sporočili pošiljata v kanal; glavna gorutina izpisuje vse, kar se pojavi v kanalu
// če zahtevamo več podatkov, kot jih gorutine pošljejo v kanal, pride do smrtnega objema
// opazujemo, kdaj gorutine zaključijo

package main

import (
	"fmt"
	"strconv"
)

const printouts = 10

var stringStream chan string

func hello(s string) {
	fmt.Println("Goroutine", s, "start")
	for i := 0; i < printouts; i++ {
		stringStream <- s + "-" + strconv.Itoa(i)
		fmt.Println("Iteration", i, "done")
	}
	fmt.Println("Goroutine", s, "done")
}

func main() {

	stringStream = make(chan string)

	go hello("hello")
	go hello("world")

	var msg string
	for i := 0; i < 2*printouts; i++ {
		msg = <-stringStream
		fmt.Print(msg, " ")
	}
	fmt.Println()
}
