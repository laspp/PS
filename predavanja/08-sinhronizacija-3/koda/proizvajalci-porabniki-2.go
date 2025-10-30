// Proizvajalec-porabnik
// rešitev s kanali
// 		kanal s pogojnimi spremenljivkami in ključavnico

package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

type productData struct {
	Id int
}

var buffer []productData
var bufferIdxPut int = 0
var bufferIdxGet int = 0
var bufferNumProducts int = 0

var wgProducer, wgConsumer sync.WaitGroup
var lockBuffer sync.Mutex
var condProducers = sync.NewCond(&lockBuffer)
var condConsumers = sync.NewCond(&lockBuffer)

func createProduct(producerId int, taskId int, product *productData) {
	product.Id = 10*producerId + taskId
	fmt.Println("P   ", producerId, product.Id)
}

func consumeProduct(consumerId int, product *productData) {
	fmt.Println("\tC   ", consumerId, product.Id)
}

func producer(id int, products int, bufferSize int) {
	var product productData
	defer wgProducer.Done()
	for i := 0; i < products; i++ {
		createProduct(id, i+1, &product)

		condProducers.L.Lock()
		for bufferNumProducts == bufferSize {
			condProducers.Wait()
		}

		fmt.Println("P->b", id, product.Id)
		buffer[bufferIdxPut] = product
		bufferIdxPut = (bufferIdxPut + 1) % bufferSize
		bufferNumProducts++

		condConsumers.Signal()
		condProducers.L.Unlock()
	}
}

func consumer(id int, bufferSize int) {
	var product productData
	defer wgConsumer.Done()
	for {
		condConsumers.L.Lock()
		for bufferNumProducts == 0 {
			condConsumers.Wait()
		}

		product = buffer[bufferIdxGet]
		bufferIdxGet = (bufferIdxGet + 1) % bufferSize
		bufferNumProducts--
		fmt.Println("\tb->C", id, product.Id)

		condProducers.Signal()
		condConsumers.L.Unlock()

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

	// pripravimo medpomnilnik
	buffer = make([]productData, *bPtr)

	// zaženemo gorutine
	wgProducer.Add(*pPtr)
	for i := 0; i < *pPtr; i++ {
		go producer(i+1, *nPtr, *bPtr)
	}
	wgConsumer.Add(*cPtr)
	for i := 0; i < *cPtr; i++ {
		go consumer(i+1, *bPtr)
	}

	// počakamo, da proizvajalci zaključijo
	wgProducer.Wait()
	time.Sleep(time.Second)
}
