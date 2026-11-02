// Proizvajalec-porabnik
// rešitev s kanali
// 		rešitev s kanalom podpira tudi predajanje izdelka iz rok v roke, brez odlaganja na polico
// 		izpisom ni popolnoma za verjeti - izpis se lahko zgodi pred vpisom v kanal

package main

import (
	"flag"
	"fmt"
	"sync"
)

type productData struct {
	Id int
}

var wgProducer, wgConsumer sync.WaitGroup
var bufferStream chan productData

func createProduct(producerId int, taskId int, product *productData) {
	product.Id = 10*producerId + taskId
	fmt.Println("P   ", producerId, product.Id)
}

func consumeProduct(consumerId int, product *productData) {
	fmt.Println("\tC   ", consumerId, product.Id)
}

func producer(id int, products int) {
	var product productData
	defer wgProducer.Done()
	for i := 0; i < products; i++ {
		createProduct(id, i+1, &product)
		// vpišemo v medpomnilnik
		fmt.Println("P->b", id, product.Id)
		bufferStream <- product
	}
}

func consumer(id int) {
	defer wgConsumer.Done()
	// iz medpomnilnika beremo dokler kanal ni zaprt in prazen
	for product := range bufferStream {
		fmt.Println("\tb->C", id, product.Id)
		consumeProduct(id, &product)
	}
}

func main() {
	// preberemo argumente
	pPtr := flag.Int("p", 1, "# of producers")
	cPtr := flag.Int("c", 1, "# of consumers")
	bPtr := flag.Int("b", 1, "buffer size")
	nPtr := flag.Int("n", 5, "number of products per producer")
	flag.Parse()

	// pripravimo kanal
	bufferStream = make(chan productData, *bPtr)

	// zaženemo gorutine
	wgProducer.Add(*pPtr)
	for i := 0; i < *pPtr; i++ {
		go producer(i+1, *nPtr)
	}
	wgConsumer.Add(*cPtr)
	for i := 0; i < *cPtr; i++ {
		go consumer(i + 1)
	}

	// počakamo, da proizvajalci zaključijo
	wgProducer.Wait()

	// zapremo kanal
	close(bufferStream)

	// počakamo, da porabniki izpraznijo kanal in zaključijo
	wgConsumer.Wait()
}
