// Problem pisateljev in bralcev
// rešitev s semaforjema, ki ju izvedemo s kanaloma

package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var semBook chan struct{}               	// števni semafor z vrednostmi med 0 in *readersPtr
var semGatekeeper = make(chan struct{}, 1)	// binarni semafor, ki preprečuje hkratno delo s kanalom semBook

func writer(id int, cycles int, places int) {
	defer wg.Done()

	for i := 0; i < cycles; i++ {
		// pisatelj lahko piše, ko zasede vsa mesta v kanalu
		// semGatekeeper prepreči, da bi hkrati vstopila dva pisatelja
		semGatekeeper <- struct{}{}
		for k := 0; k < places; k++ {
			semBook <- struct{}{}
		}
		<-semGatekeeper

		fmt.Println("Writer", id, "start", i)
		time.Sleep(time.Duration(id) * time.Millisecond)
		fmt.Println("Writer", id, "finish", i)

		for k := 0; k < places; k++ {
			<-semBook
		}

		time.Sleep(time.Duration(id) * time.Millisecond)
	}
}

func reader(id int) {

	for {
		// med tem, ko pisatelj pridobiva mesta v kanalu, morajo bralci čakati
		semGatekeeper <- struct{}{}
		semBook <- struct{}{}
		<-semGatekeeper

		fmt.Println("Reader", id, "start")
		time.Sleep(time.Duration(id) * time.Millisecond)
		fmt.Println("Reader", id, "finish")

		<-semBook

		time.Sleep(time.Duration(id) * time.Millisecond)
	}
}

func main() {
	// preberemo argumente
	writersPtr := flag.Int("w", 2, "# of writers")
	readersPtr := flag.Int("r", 4, "# of readers")
	cyclesPtr := flag.Int("c", 10, "# of cycles")
	flag.Parse()

	// kapaciteta kanala mora omogočati branje vseh bralcev naenkrat
	semBook = make(chan struct{}, *readersPtr)

	// zaženemo pisatelje
	wg.Add(*writersPtr)
	for i := 1; i <= *writersPtr; i++ {
		go writer(i, *cyclesPtr, *readersPtr)
	}
	// zaženemo bralce
	for i := 1; i <= *readersPtr; i++ {
		go reader(i)
	}
	// počakamo, da pisatelji zaključijo
	wg.Wait()
}
