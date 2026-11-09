/*
Primer proizvajalec porabnik
Ustvarimo proizvajalce in porabnike, na koncu preverimo, kaj se je zgodilo z zahtevami.
Namesto lastnega kanala quit za zaustavitev uporabimo kontekst (context.Context).
Kontekst s časovno omejitvijo se po dveh sekundah sam prekliče.
Najprej ustavimo proizvajalce, in nato počakamo na porabnike.
*/
package main

import (
	"context"
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

func (p *producer) start(ctx context.Context, data chan<- int) {
	for {
		select {
		// Kanal ctx.Done() se zapre, ko je kontekst preklican
		case <-ctx.Done():
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
	// Ustvarimo kontekst, ki se po dveh sekundah samodejno prekliče.
	// Funkcijo cancel poklicemo, da sprostimo vire konteksta,
	// z njo pa bi lahko kontekst preklicali tudi prej.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	data := make(chan int)
	nProducers := 10
	nConsumers := 10
	var wgProducer, wgConsumer sync.WaitGroup

	producers := make([]producer, nProducers)
	consumers := make([]consumer, nConsumers)

	for i := range producers {
		producers[i].id = i
		wgProducer.Go(func() { producers[i].start(ctx, data) })
	}

	for i := range consumers {
		consumers[i].id = i
		wgConsumer.Go(func() { consumers[i].start(data) })
	}

	// Počakamo, da se kontekst prekliče, in izpišemo razlog
	<-ctx.Done()
	fmt.Println("Zaustavitev:", ctx.Err())
	wgProducer.Wait()
	close(data)
	wgConsumer.Wait()
	checkWork(producers, consumers)
}
