/*
Primer proizvajalec porabnik
Ustvarimo proizvajalce in porabnike, na koncu preverimo, kaj se je zgodilo z zahtevami.
Težava, ne počakamo, da proizvajalci in porabniki zaključijo
*/
package main

import (
	"fmt"
	"time"
)

type producer struct {
	id      int
	counter int
}

type consumer struct {
	id      int
	counter int
}

func (p *producer) start(data chan<- int) {
	for {
		p.counter++
		data <- p.counter
	}
}

func (c *consumer) start(data <-chan int) {
	for {
		<-data
		c.counter++
	}
}

func checkWork(ps []producer, cs []consumer) {
	total := 0
	for _, p := range ps {
		total += p.counter
		fmt.Println("P", p.id, p.counter)
	}

	for _, c := range cs {
		total -= c.counter
		fmt.Println("C", c.id, c.counter)
	}
	// TOTAL je razlika med številom proizvedenih in porabljenih podatkov.
	fmt.Println("TOTAL:", total)
}

func main() {
	data := make(chan int)
	nProducers := 10
	nConsumers := 10

	producers := make([]producer, nProducers)
	consumers := make([]consumer, nConsumers)

	for i := range producers {
		producers[i].id = i
		go producers[i].start(data)
	}

	for i := range consumers {
		consumers[i].id = i
		go consumers[i].start(data)
	}

	time.Sleep(2 * time.Second)
	checkWork(producers, consumers)
}
