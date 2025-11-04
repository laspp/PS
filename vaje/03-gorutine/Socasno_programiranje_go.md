# Sočasno programiranje v jeziku go

Jezik **go** je bil zasnovan z mislijo na sočasno programiranje in kot tak nudi ustrezno podporo skozi vgrajene programske stavke in standardno knjižnico jezika. Osnovni gradnik so gorutine, ki služijo kot abstrakcija niti operacijskega sistema. Pri čemer ne velja, da je 1 gorutina enaka 1 programski niti. Gorutine so veliko lažje od običajnih niti, saj zahtevajo manj režijskih stroškov pri upravljanju. To pomeni, da jih lahko zaganjamo v velikem številu in s tem ne ohromimo sistema. 

Gorutino ustvarimo s pomočjo programskega stavka `go`.
Primer ustvarjanja gorutin:
```Go
package main

import (
	"fmt"
	"time"
)

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}

func main() {
	// Zaženemo funkcijo f znotraj glavne gorutine
	f("Glavna nit")

	// Ustvarimo novo gorutino v kateri se izvaja funkcija f
	go f("Prva gorutina")

	// Ustvarimo novo gorutino, ki izvede anonimno funkcijo
	go func(msg string) {
		fmt.Println(msg)
	}("Druga gorutina")

	// Počakamo, da se vse gorutine zaključijo
	time.Sleep(time.Second)
	fmt.Println("KONEC")
}
```
[gorutine.go](koda/gorutine.go)

Na koncu programa moramo uporabiti funkcijo `time.Sleep(time.Second)`, zato da počakamo, da se vse gorutine zaključijo. V nasprotnem primeru bi se glavna gorutina lahko zaključila pred ostalimi in bi dobili nepopolni izpis. Seveda, to ni ravno najboljši način izvedbe sinhronizacije niti.

## Kanali

Kanali služijo sinhronizaciji in komunikaciji med gorutinami. Kanal definiramo s pomočjo ključne besede `chan`. Kanalu pripišemo tudi podatkovni tip. S tem povemo kakšen tip podatkov lahko po njem prenašamo.
```Go
    // Kanal, ki sprejema cela števila
    var ch chan int
    ch = make(chan int)
```
Pri delu s kanali uporabljamo operator `<-` za **pisanje** in **branje** iz kanala.
```Go
    // Pisanje v kanal
    ch <- 42
    // Branje iz kanala
    v := <- ch
```
Kanalu lahko določimo kapaciteto, privzeta kapaciteta je 0. Če je kanal poln (kanal s kapaciteto 0 je vedno poln) potem tak kanal blokira izvajanje gorutine, ko le-ta izvede branje ali pisanje v kanal. Izvajanje gorutine se nadaljuje, ko neka druga gorutina pošlje ali prebere vrednost iz kanala.

```Go
    // Ustvarimo kanal s kapaciteto 2
    var ch := make(chan int, 2)
    // V kanal lahko zapišemo dve vrednosti preden ta blokira goruitno
    kanal <- 1
    kanal <- 2
    // tukaj se bo izvajanje ustavilo, dokler neka druga gorutine ne prebere vrednosti iz kanala
    kanal <- 3
```

S kanali lahko enostavno rešimo problem čakanja na končanje zagnanih gorutin. Primer izvedbe sinhronizacije gorutin s pomočjo kanalov:
```Go
/*
Program kanali demonstrira uporabo kanalov v programskem jeziku go
*/
package main

import (
	"fmt"
	"time"
)

func worker(id int, done chan bool) {
	fmt.Println(id, "Delam ...")
	time.Sleep(time.Second)
	fmt.Println(id, "Zaključil")
	done <- true
}

func main() {
	// Ustvarimo kanal s kapaciteto 3
	workers := 3
	done := make(chan bool, workers)
	// Zaženemo delavce
	for w := 0; w < workers; w++ {
		go worker(w, done)
	}
	// Počakamo, da delavci zaključijo
	for w := 0; w < workers; w++ {
		<-done
	}
}
```
[kanali.go](koda/kanali.go)

### Stavek select
Pri delu z gorutinami in kanali nam **go** nudi stavek `select`, ki ima podobno zgradbo kot stavek `switch`, vendar služi spremljanju dogajanja na več kanalih hkrati. Stavek `select` blokira izvajanje gorutine dokler se ne zgodi dogodek na enem od kanalov, ki jih spremlja. Takrat izbere vejo, ki se je sprožila, in jo izvede. V primeru, da je hkrati pripravljenih več vej, izbere naključno vejo. Uporabimo lahko tudi privzeto vejo `default`, ki se izvede, ko nobena druga veja ni pripravljena.

Primer uporabe stavka `select`, kjer čakamo na pritisk gumba `Enter`. Dokler se to ne zgodi pa zaganjamo nove gorutine `worker`:
```Go
/*
Program select prikazuje uporabo stavka select pri delu s kanali v programskem jeziku go
*/
package main

import (
	"fmt"
	"time"
)

// Funkcija, ki čaka na pritisk tipke "Enter"
func readKey(input chan bool) {
	fmt.Scanln()
	input <- true
}

// Delavec
func worker(id int, done chan bool) {
	fmt.Print("Delavec ", id)
	time.Sleep(2 * time.Second)
	fmt.Print("Končal")
	done <- true
}

func main() {
	input := make(chan bool)
	done := make(chan bool)
	w := 0
	// Zaženemo gorutino, ki čaka na pritisk tipke
	go readKey(input)
	// Zaženemo prvega delavca
	go worker(w, done)
	// Anonimna funkcija z neskončno zanko
	func() {
		for {
			select {
			// Pritisk tipke: končamo
			case <-input:
				return
			// Delavec zaključil, zaženimo novega
			case <-done:
				fmt.Println()
				w = w + 1
				go worker(w, done)
			// Če se nič ne zgodi izvedemo privzeto akcijo
			default:
				time.Sleep(200 * time.Millisecond)
				fmt.Print(".")
			}
		}
	}()
	// Počakamo na zadnjega delavca
	<-done
}
```
[select.go](koda/select.go)
## Konstrukti za sinhronizacijo
Poleg kanalov, imamo pri delu z gorutinami na voljo še vgrajen paket `sync`, v katerem najdemo konstrukte, kot so ključavnice (angl. Mutex), in čakalne skupine (angl. WaitGroups). Primer uporabe omenjenih konstruktov:
```Go
/*
Program sync prikazuje uporabo sinhronizacijskih konstruktov v paketu sync
*/
package main

import (
	"fmt"
	"sync"
)

// Definiramo čakalno skupino
var wg sync.WaitGroup
var wc int

// Definiramo ključavnico
var lock sync.Mutex

// Delavec, ki povečuje števec
func workerInc(id int) {
	defer wg.Done()
	lock.Lock()
	wc++
	lock.Unlock()
}

// delavec, ki zmanjšuje števec
func workerDec(id int) {
	defer wg.Done()
	lock.Lock()
	wc--
	lock.Unlock()
}

func main() {
	workers := 100
	// Čakalno skupino inicializiramo z želenim številom delavcev
	wg.Add(2 * workers)
	// Zaženemo delavce
	for i := 0; i < workers; i++ {
		go workerInc(i)
		go workerDec(i)
	}
	// Počakamo, da delavci zaključijo
	wg.Wait()
	// Izpišemo končni rezultat
	// Kaj se zgodi, če iz delavcev odstranimo zaklepanje in odklepanje ključavnic?
	fmt.Println("Števec: ", wc)
}
```
[sync.go](koda/sync.go)

# Domača naloga 2

Napišite program v jeziku Go, ki bo omogočal spremljanje in izpis meritev senzorjev temperature, vlage in tlaka s pomočjo gorutin. Rešitev oddajte preko [spletne učilnice](https://ucilnica.fri.uni-lj.si/mod/assign/view.php?id=53926).

**Navodila:**

Definirajte strukturo `Meritev`:
```Go
type Meritev struct {
	vrsta string
    vrednost float32
}
```
Ustvarite 3 gorutine, ki bodo simulirale senzorje temperature, vlage in tlaka in vsakih 100 ms pošiljale meritve v kanal s pomočjo strukture `Meritev`. Niz `vrsta` naj določa tip meritve; torej `temperatura`, `vlaga` ali `tlak`. Vrednosti generirajte naključno na intervalu, ki je smiseln za dano količino.

Glavna gorutina naj s pomočjo stavka `select` bere meritve iz kanala in jih izpisuje na zaslon. Kanal za sporočanje meritev naj ima kapaciteto 3. Glavna gorutina naj preverja tudi ali je merilni sistem postal neodziven (5 sekund ni bila preko kanala prejeta nobena meritev od katerega koli senzorja). Če ugotovi, da je sistem neodziven naj izpiše obvestilo na zaslon in zaključi. Prav tako naj aplikacija omogoča, da jo s pritiskom na tipko `Enter` ustavimo. Za slednje uporabite dodatno gorutino, ki čaka na pritisk tipke in o tem obvesti glavno gorutino preko ločenega kanala.  

V pomoč vam bodo naslednje funkcije:
```Go
fmt.Scanln()
time.After()
time.Sleep()
rand.Float32()
```

**Rok za oddajo: 11. 11. 2025**

