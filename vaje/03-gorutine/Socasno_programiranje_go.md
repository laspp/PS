# Sočasno programiranje v jeziku Go

Jezik Go je bil zasnovan z mislijo na sočasno programiranje in kot tak nudi ustrezno podporo z vgrajenimi programskimi stavki in standardno knjižnico jezika. Osnovni gradnik so gorutine, ki služijo kot abstrakcija niti operacijskega sistema, pri čemer ne velja, da je ena gorutina enaka eni niti operacijskega sistema. Gorutine so veliko lažje od običajnih niti, saj zahtevajo manj režijskih stroškov pri upravljanju.

Za izvajanje gorutin skrbi razvrščevalnik v izvajalnem okolju (angl. *runtime*) jezika Go, ki $M$ gorutin razporeja na $N$ niti operacijskega sistema. Koliko niti lahko hkrati izvaja kodo Go, določa nastavitev `GOMAXPROCS`, ki je privzeto enaka številu logičnih jeder procesorja. Ko gorutina blokira (npr. čaka na kanal), razvrščevalnik na isti niti zažene drugo gorutino. Sklad gorutine na začetku obsega le nekaj KB in se po potrebi samodejno povečuje, medtem ko ima nit operacijskega sistema običajno fiksen sklad velikosti nekaj MB. Zato lahko brez težav hkrati poganjamo tudi več deset tisoč gorutin.

Primeri iz tega poglavja so urejeni v en modul, pri čemer ima vsak program svoj direktorij. Najdete jih v mapi [koda](koda/):
```
koda/
├── go.mod              
├── 01-gorutine/
│   └── gorutine.go      
├── 02-kanali/
│   └── kanali.go
├── 03-select/
│   └── select.go
└── 04-sync/
    └── sync.go
```
Posamezen primer zaženemo iz mape `koda` z ukazom `go run <ime_mape_primera>`.

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
	f("Glavna gorutina")

	// Ustvarimo novo gorutino, v kateri se izvaja funkcija f
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
[gorutine.go](koda/01-gorutine/gorutine.go)

Na koncu programa moramo uporabiti funkcijo `time.Sleep(time.Second)`, da počakamo, da se vse gorutine zaključijo. V nasprotnem primeru bi se glavna gorutina lahko zaključila pred ostalimi in bi dobili nepopoln izpis. Seveda to ni ravno najboljši način izvedbe sinhronizacije niti.

## Kanali

Kanali služijo sinhronizaciji in komunikaciji med gorutinami. Kanal definiramo s pomočjo ključne besede `chan`. Kanalu pripišemo tudi podatkovni tip. S tem povemo, kakšen tip podatkov lahko po njem prenašamo.
```Go
    // Kanal, ki sprejema cela števila
    var ch chan int
    ch = make(chan int)
```
Pri delu s kanali uporabljamo operator `<-` za **pisanje** v kanal in **branje** iz kanala.
```Go
    // Pisanje v kanal
    ch <- 42
    // Branje iz kanala
    v := <-ch
```
Kanalu lahko določimo kapaciteto (privzeta je 0). Pisanje v kanal blokira gorutino, dokler je medpomnilnik kanala poln, branje pa, dokler je prazen. Pri kanalu brez medpomnilnika (kapaciteta 0) se morata pošiljatelj in prejemnik "srečati": pisanje blokira, dokler druga gorutina ne prebere vrednosti, in obratno.

```Go
    // Ustvarimo kanal s kapaciteto 2
    ch := make(chan int, 2)
    // V kanal lahko zapišemo dve vrednosti, preden ta blokira gorutino
    ch <- 1
    ch <- 2
    // Tukaj se bo izvajanje ustavilo, dokler neka druga gorutina ne prebere vrednosti iz kanala
    ch <- 3
```

S kanali lahko enostavno rešimo problem čakanja na končanje zagnanih gorutin. Primer izvedbe sinhronizacije gorutin s pomočjo kanalov:
```Go
/*
Program kanali demonstrira uporabo kanalov v programskem jeziku Go
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
[kanali.go](koda/02-kanali/kanali.go)

### Zapiranje kanalov
Pošiljatelj lahko kanal zapre s funkcijo `close`, s čimer prejemnikom sporoči, da novih vrednosti ne bo več. Branje iz zaprtega kanala nikoli ne blokira: najprej vrne vrednosti, ki so še ostale v medpomnilniku, nato pa ničelno vrednost podatkovnega tipa kanala. Ali je bila vrednost res poslana ali je kanal zaprt, preverimo z drugo povratno vrednostjo pri branju. Pisanje v zaprt kanal povzroči paniko (angl. *panic*). Z zanko `for range` beremo vrednosti iz kanala, dokler ta ni zaprt in prazen.
```Go
    ch := make(chan int, 3)
    go func() {
        for i := 0; i < 3; i++ {
            ch <- i
        }
        // Sporočimo, da novih vrednosti ne bo več
        close(ch)
    }()
    // Zanka se zaključi, ko je kanal zaprt in prazen
    for v := range ch {
        fmt.Println(v)
    }
    // Branje iz zaprtega kanala vrne ničelno vrednost in ok == false
    v, ok := <-ch
    fmt.Println(v, ok) // 0 false
```
Kanal naj vedno zapre pošiljatelj, nikoli prejemnik, saj le pošiljatelj ve, kdaj novih vrednosti ne bo več.

### Usmerjeni kanali
Pri podajanju kanala funkciji lahko omejimo smer prenosa podatkov. Kanal tipa `chan<- T` dovoljuje samo pisanje, kanal tipa `<-chan T` pa samo branje. Dvosmerni kanal se ob klicu funkcije samodejno pretvori v usmerjenega. Prevajalnik nato javi napako, če funkcija kanal uporabi v napačni smeri. Tako že iz glave funkcije vidimo njeno vlogo, hkrati pa se izognemo nekaterim napakam.
```Go
// Funkcija lahko v kanal samo piše
func writer(out chan<- int) {
    out <- 42
}

// Funkcija lahko iz kanala samo bere
func reader(in <-chan int) {
    fmt.Println(<-in)
}
```

### Stavek select
Pri delu z gorutinami in kanali nam Go nudi stavek `select`, ki ima podobno zgradbo kot stavek `switch`, vendar služi spremljanju dogajanja na več kanalih hkrati. Stavek `select` blokira izvajanje gorutine, dokler se ne zgodi dogodek na enem od kanalov, ki jih spremlja. Takrat izbere vejo, ki se je sprožila, in jo izvede. V primeru, da je hkrati pripravljenih več vej, izbere naključno vejo. Uporabimo lahko tudi privzeto vejo `default`, ki se izvede, ko nobena druga veja ni pripravljena.

Primer uporabe stavka `select`, kjer čakamo na pritisk tipke `Enter`. Dokler se to ne zgodi, zaganjamo nove gorutine `worker`:
```Go
/*
Program select prikazuje uporabo stavka select pri delu s kanali v programskem jeziku Go
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
	fmt.Print(" Končal")
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
			// Če se nič ne zgodi, izvedemo privzeto akcijo
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
[select.go](koda/03-select/select.go)

## Konstrukti za sinhronizacijo
Poleg kanalov imamo pri delu z gorutinami na voljo še paket `sync` iz standardne knjižnice, v katerem najdemo konstrukte, kot so ključavnice (angl. *mutex*) in čakalne skupine (angl. *wait group*). Primer uporabe omenjenih konstruktov:
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
func workerInc() {
	defer wg.Done()
	lock.Lock()
	wc++
	lock.Unlock()
}

// Delavec, ki zmanjšuje števec
func workerDec() {
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
		go workerInc()
		go workerDec()
	}
	// Počakamo, da delavci zaključijo
	wg.Wait()
	// Izpišemo končni rezultat
	// Kaj se zgodi, če iz delavcev odstranimo zaklepanje in odklepanje ključavnic?
	fmt.Println("Števec:", wc)
}
```
[sync.go](koda/04-sync/sync.go)

Od različice Go 1.25 dalje ima čakalna skupina tudi metodo `Go`, ki v novi gorutini zažene podano funkcijo ter namesto nas pokliče `Add(1)` pred zagonom in `Done()` ob zaključku funkcije. Zanko iz zgornjega primera bi lahko zapisali takole (funkciji delavcev v tem primeru ne kličeta `wg.Done()`):
```Go
	for i := 0; i < workers; i++ {
		wg.Go(workerInc)
		wg.Go(workerDec)
	}
	wg.Wait()
```

Če iz delavcev odstranimo ključavnico, pride do tveganega stanja (angl. *race condition*): več gorutin hkrati bere in piše spremenljivko `wc`, zato končna vrednost števca ni nujno enaka 0. Takšne napake se pogosto pokažejo le občasno, zato jih je težko odkriti s samim poganjanjem programa. Pri tem nam pomaga detektor tveganih stanj, ki je vgrajen v orodja jezika Go. Vklopimo ga z zastavico `-race`:
```Bash
go run -race ./04-sync
```
Program se prevede z dodatno kodo, ki med izvajanjem spremlja dostope do pomnilnika. Ko zazna, da dve gorutini brez sinhronizacije dostopata do iste spremenljivke in je vsaj eden od dostopov pisanje, izpiše opozorilo `WARNING: DATA RACE` skupaj z mestoma v kodi, kjer sta se dostopa zgodila. Detektor najde le tvegana stanja, ki se med danim izvajanjem dejansko zgodijo, program pa z njim teče občutno počasneje in porabi več pomnilnika. Zato ga uporabljamo med razvojem in testiranjem. 


# Domača naloga 2

Napišite program v jeziku Go, ki bo omogočal spremljanje in izpis meritev senzorjev temperature, vlage in tlaka s pomočjo gorutin. Rešitev oddajte preko [spletne učilnice](https://ucilnica.fri.uni-lj.si/mod/assign/view.php?id=53926).

**Navodila:**

Definirajte strukturo `Meritev`:
```Go
type Meritev struct {
	vrsta    string
	vrednost float32
}
```
Ustvarite 3 gorutine, ki bodo simulirale senzorje temperature, vlage in tlaka ter vsakih 500 ms pošiljale meritve v kanal s pomočjo strukture `Meritev`. Vsi trije senzorji naj pošiljajo meritve v isti kanal. Niz `vrsta` naj določa tip meritve, torej `temperatura`, `vlaga` ali `tlak`. Vrednosti generirajte naključno iz intervala, ki je smiseln za dano količino (npr. temperatura od -20 do 40 °C, vlaga od 0 do 100 % in tlak od 950 do 1050 hPa). Za generiranje naključnih števil lahko uporabite paket `math/rand` ali `math/rand/v2`.

Glavna gorutina naj s pomočjo stavka `select` bere meritve iz kanala in jih izpisuje na zaslon, vsako v svojo vrstico in z enoto, na primer:
```
temperatura: 21.37 °C
vlaga: 45.02 %
tlak: 1013.25 hPa
```
Kanal za sporočanje meritev naj ima kapaciteto 3. Glavna gorutina naj preverja tudi, ali je merilni sistem postal neodziven (5 sekund ni prek kanala prejela nobene meritve od nobenega senzorja). Če ugotovi, da je sistem neodziven, naj izpiše obvestilo na zaslon in zaključi. Da bo mogoče preveriti zaznavanje neodzivnosti, naj vsak senzor po naključnem času (npr. med 5 in 20 sekundami od zagona) preneha pošiljati meritve, kot da je prišlo do okvare. Prav tako naj aplikacija omogoča, da jo s pritiskom na tipko `Enter` ustavimo. Za slednje uporabite dodatno gorutino, ki čaka na pritisk tipke in o tem obvesti glavno gorutino prek ločenega kanala.

V pomoč vam bodo naslednje funkcije:
```Go
fmt.Scanln()
time.After()
time.Sleep()
rand.Float32()
```

**Rok za oddajo: 2. 11. 2026**

