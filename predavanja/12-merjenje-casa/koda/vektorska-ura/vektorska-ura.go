// Komunikacija po protokolu TCP z beleženjem z vektorsko uro
//		z beleženjem je dopolnjen primer tcp-nizi
//
// 		odjemalec pošlje strežniku sporočilo, dopolnjeno s časovnim žigom
//		strežnik sporočilo dopolni, zamenja časovni žig in ga pošlje odjemalcu
//		pred pošiljanjem počaka 5 s, da lažje pokažemo hkratno streženje več odjemalcem
//
//		med izvajanjem beležimo lokalne dogodke in prenos podatkov
//
// zaženemo strežnik
// 		go run *.go
// zaženemo enega ali več odjemalcev
//		go run *.go -s [ime strežnika] -p [vrata] -m [sporočilo]
// za [ime strežnika] in [vrata] vpišemo vrednosti, ki jih izpiše strežnik ob zagonu
//
// za vizualizacijo najprej naredimo združeni log
// 		~/go/bin/GoVector --log_type shiviz --log_dir . --outfile Log-Visu.log
// nato datoteko Log-Visu.log odpremo v spletni aplikaciji
// 		https://bestchai.bitbucket.io/shiviz/

package main

import (
	"flag"
	"fmt"
)

func main() {
	// preberemo argumente iz ukazne vrstice
	sPtr := flag.String("s", "", "server URL")
	pPtr := flag.Int("p", 9876, "port number")
	mStr := flag.String("m", "world", "message")
	flag.Parse()

	// zaženemo strežnik ali odjemalca
	url := fmt.Sprintf("%v:%v", *sPtr, *pPtr)
	if *sPtr == "" {
		Server(url)
	} else {
		Client(url, *mStr)
	}
}
