// Komunikacija po protokolu RPC
//
// 		strežnik ustvari in vzdržuje shrambo nalog TodoStorage
//		odjemalec nad shrambo izvaja operacije CRUD
//		dve verziji: HTTP in TCP, z verzijo HTTP lahko poskrbimo tudi za avtorizacijo odjemalca na strežniku
//
// zaženemo strežnik
// 		go run *.go -c [tip povezave http ali tcp]
// zaženemo enega ali več odjemalcev
//		go run *.go -s [ime strežnika] -p [vrata] -c [tip povezave http ali tcp]
// za [ime strežnika] in [vrata] vpišemo vrednosti, ki jih izpiše strežnik ob zagonu
// strežnik in odjemalec morata uporabljati isti tip povezave
//
// pri uporabi SLURMa lahko s stikalom --nodelist=[vozlišče] določimo vozlišče, kjer naj se program zažene

package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	// preberemo argumente iz ukazne vrstice
	sPtr := flag.String("s", "", "server URL")
	pPtr := flag.Int("p", 9876, "port number")
	cStr := flag.String("c", "http", "connection type")
	flag.Parse()

	// zaženemo strežnik ali odjemalca
	url := fmt.Sprintf("%v:%v", *sPtr, *pPtr)
	connHTTP := strings.HasPrefix(strings.ToUpper(*cStr), "H")
	if *sPtr == "" {
		Server(url, connHTTP)
	} else {
		Client(url, connHTTP)
	}
}
