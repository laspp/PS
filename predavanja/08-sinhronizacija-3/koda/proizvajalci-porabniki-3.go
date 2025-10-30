// Proizvajalec-porabnik
// rešitev s kanali
// 		kanal je za ena manjši od velikosti medpomnilnika, da je rešitev enakovredna rešitvi s pogojnimi spremenljivkami
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
	// iz medpomnilnika beremo dokler ni prazen (kanal zaprt)
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
	bufferStream = make(chan productData, *bPtr-1)

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

	// počakamo, do porabniki izpraznijo kanal in zaključijo
	wgConsumer.Wait()
}
