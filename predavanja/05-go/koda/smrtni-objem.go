package main

import (
	"fmt"
)

var dataStream = make(chan string)

func writer() {
	dataStream <- "13" // gorutina nikoli ne zapiše vrednosti v kanal
	return
}

func main() {
	go writer()
	value := <-dataStream // glavna gorutina blokira ob čakanju na podatek --> smrtni objem
	fmt.Println("Value:", value)
}
