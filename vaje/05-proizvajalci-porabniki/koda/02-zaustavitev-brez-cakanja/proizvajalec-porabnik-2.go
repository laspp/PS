/*
Primer proizvajalec porabnik
Ustvarimo proizvajalce in porabnike, na koncu preverimo, kaj se je zgodilo z zahtevami.
Glavna nit pošlje signal za zaustavitev. Težava, ne vemo kdaj se gorutine zares ustavijo.
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

func (p *producer) start(data chan<- int, quit <-chan struct{}) {
	for {
		select {
		case <-quit:
			return
		default:
			p.counter++
			data <- p.counter
		}
	}
}

func (c *consumer) start(data <-chan int, quit <-chan struct{}) {
	for {
		select {
		case <-data:
			c.counter++
		case <-quit:
			return
		}
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
	quit := make(chan struct{})
	nProducers := 10
	nConsumers := 10

	producers := make([]producer, nProducers)
	consumers := make([]consumer, nConsumers)

	for i := range producers {
		producers[i].id = i
		go producers[i].start(data, quit)
	}

	for i := range consumers {
		consumers[i].id = i
		go consumers[i].start(data, quit)
	}

	time.Sleep(2 * time.Second)
	close(quit)
	checkWork(producers, consumers)
}
