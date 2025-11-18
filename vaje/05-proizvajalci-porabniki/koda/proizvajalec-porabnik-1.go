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

func (p *producer) start(id int, interval time.Duration, data chan<- int) {
	p.id = id
	for {
		p.counter++
		data <- p.counter
		time.Sleep(interval)
	}
}

func (c *consumer) start(id int, interval time.Duration, data <-chan int) {
	c.id = id
	for {
		<-data
		c.counter++
		time.Sleep(interval)
	}
}

func checkWork(ps []producer, cs []consumer) {
	total := 0
	for i, p := range ps {
		total += p.counter
		fmt.Println("P", i, p.counter)
	}

	for i, c := range cs {
		total -= c.counter
		fmt.Println("C", i, c.counter)
	}
	fmt.Println("TOTAL:", total)
}

func main() {
	data := make(chan int)
	nProducers := 10
	nConsumers := 10

	producers := make([]producer, nProducers)
	consumers := make([]consumer, nConsumers)
	intervalProducers := 10 * time.Millisecond
	intervalConsumers := 10 * time.Millisecond

	for i := range producers {
		go producers[i].start(i, intervalProducers, data)
	}

	for i := range consumers {
		go consumers[i].start(i, intervalConsumers, data)
	}

	time.Sleep(2 * time.Second)
	checkWork(producers, consumers)

}
