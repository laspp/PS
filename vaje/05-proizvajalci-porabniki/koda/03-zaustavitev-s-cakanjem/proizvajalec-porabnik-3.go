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

func (c *consumer) start(data <-chan int) {
	// Zanka se zaključi, ko je kanal zaprt in prazen
	for range data {
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
	quit := make(chan struct{})
	nProducers := 10
	nConsumers := 10
	var wgProducer, wgConsumer sync.WaitGroup

	producers := make([]producer, nProducers)
	consumers := make([]consumer, nConsumers)

	// Metoda Go zažene funkcijo v novi gorutini ter namesto nas
	// pokliče Add(1) pred zagonom in Done() ob zaključku
	for i := range producers {
		producers[i].id = i
		wgProducer.Go(func() { producers[i].start(data, quit) })
	}

	for i := range consumers {
		consumers[i].id = i
		wgConsumer.Go(func() { consumers[i].start(data) })
	}

	time.Sleep(2 * time.Second)
	close(quit)
	wgProducer.Wait()
	close(data)
	wgConsumer.Wait()
	checkWork(producers, consumers)
}
