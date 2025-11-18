/*
Primer proizvajalec porabnik
Ustvarimo proizvajalce in porabnike, na koncu preverimo, kaj se je zgodilo z zahtevami.
Glavna nit pošlje signal za zaustavitev. Najprej ustavimo proizvajalce, in nato počakamo na porabnike.
*/
package main

import (
	"fmt"
	"sync"
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

var wgProducer, wgConsumer sync.WaitGroup

func (p *producer) start(id int, interval time.Duration, data chan<- int, quit <-chan struct{}) {
	defer wgProducer.Done()
	p.id = id
	for {
		select {
		case <-quit:
			return
		default:
			p.counter++
			data <- p.counter
		}
		time.Sleep(interval)
	}
}

func (c *consumer) start(id int, interval time.Duration, data <-chan int) {
	defer wgConsumer.Done()
	c.id = id
	for {
		if _, more := <-data; more {
			c.counter++
		} else {
			return
		}
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
	data := make(chan int, 1000)
	quit := make(chan struct{})
	nProducers := 10
	nConsumers := 10
	wgProducer.Add(nProducers)
	wgConsumer.Add(nConsumers)

	producers := make([]producer, nProducers)
	consumers := make([]consumer, nConsumers)
	intervalProducers := 0 * time.Millisecond
	intervalConsumers := 0 * time.Millisecond

	for i := range producers {
		go producers[i].start(i, intervalProducers, data, quit)
	}

	for i := range consumers {
		go consumers[i].start(i, intervalConsumers, data)
	}

	time.Sleep(2 * time.Second)
	close(quit)
	wgProducer.Wait()
	close(data)
	wgConsumer.Wait()
	checkWork(producers, consumers)

}
