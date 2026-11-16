// Komunikacija po protokolu gRPC
//
// 		strežnik ustvari in vzdržuje shrambo nalog TodoStorage
//		odjemalec nad shrambo izvaja operacije CRUD
//		datoteka protobufStorage/protobufStorage.proto podaja strukturo sporočil gRPC
//			prevajalnik protoc iz datoteke *.proto ustvari strukture in metode, ki jih uporabljamo pri klicanju oddaljenih metod
//			navodila za prevajanje so v datoteki *.proto
//
// zaženemo strežnik
// 		go run *.go
// zaženemo enega ali več odjemalcev
//		go run *.go -s [ime strežnika] -p [vrata]
// za [ime strežnika] in [vrata] vpišemo vrednosti, ki jih izpiše strežnik ob zagonu
//
// pri uporabi SLURMa lahko s stikalom --nodelist=[vozlišče] določimo vozlišče, kjer naj se program zažene

package main

import (
	"flag"
	"fmt"
)

func main() {
	// preberemo argumente iz ukazne vrstice
	sPtr := flag.String("s", "", "server URL")
	pPtr := flag.Int("p", 9876, "port number")
	flag.Parse()

	// zaženemo strežnik ali odjemalca
	url := fmt.Sprintf("%v:%v", *sPtr, *pPtr)
	if *sPtr == "" {
		Server(url)
	} else {
		Client(url)
	}
}
