# Metode in vmesniki

Primeri iz tega poglavja so urejeni v en modul, pri čemer ima vsak program svoj direktorij. Najdete jih v mapi [koda](koda/):
```
koda/
├── go.mod             
├── 01-metode/
│   └── metode.go         
├── 02-vmesniki/
│   └── vmesniki.go
└── 03-prazen-vmesnik/
    └── prazen-vmesnik.go
```
Posamezen primer zaženemo iz mape `koda` z ukazom `go run <ime_mape_primera>`.

## Metode

Go ni zasnovan kot objektno usmerjen jezik, vsebuje pa nekatere konstrukte objektno usmerjenih jezikov. Eden od teh konstruktov so metode.

Metoda je funkcija, ki vsebuje dodaten argument, to je prejemnik (angl. *receiver*). Podatkovni tip prejemnika je lahko struktura, možni pa so tudi drugi podatkovni tipi, ki jih definiramo znotraj paketa. Podatkovni tipi, definirani izven trenutnega paketa, niso dovoljeni. Sintaksa metode je naslednja:

```Go
func (receiver type) methodName(parameter list) returnType {
}
```

Primer uporabe metod v Go:

```Go
/*
Program prikazuje uporabo metod v programskem jeziku Go
*/
package main

import (
	"fmt"
)

// Definiramo strukturo Student
type Student struct {
	name    string
	surname string
	id      string
	year    int
}

// Funkcija za izpis imena in priimka študenta
func displayNameSurname(s Student) {
	fmt.Printf("Ime: %s, Priimek: %s\n", s.name, s.surname)
}

// Metoda za izpis imena in priimka študenta
func (s Student) displayNameSurname() {
	fmt.Printf("Ime: %s, Priimek: %s\n", s.name, s.surname)
}

// Metoda, ki vrača letnik študija
func (s Student) getYear() int {
	return s.year
}

// Metoda, ki nastavi letnik študija
// Uporabimo kazalec za prejemnika
func (s *Student) setYear(year int) {
	s.year = year
}

// Metoda, ki poskuša nastaviti letnik študija
// Prejemnik je vrednost, zato metoda spremeni le svojo kopijo študenta
func (s Student) setYearCopy(year int) {
	s.year = year
}

func main() {

	// Ustvarimo novega študenta
	student1 := Student{name: "Janez", surname: "Novak", id: "63230000", year: 1}

	// Pokličemo funkcijo
	displayNameSurname(student1)

	// Pokličemo metode
	student1.displayNameSurname()
	fmt.Println(student1.getYear())
	student1.setYear(2)
	fmt.Println(student1.getYear())
	// Letnik ostane nespremenjen, saj metoda spremeni le kopijo
	student1.setYearCopy(3)
	fmt.Println(student1.getYear())
}
```
[metode.go](koda/01-metode/metode.go)

Primarna razloga za vpeljavo in uporabo metod sta naslednja:

- Go ni objektno usmerjen programski jezik in ne podpira razredov. Metode na tipih so torej način za dosego obnašanja, podobnega razredom. Metode omogočajo logično združevanje obnašanja, povezanega z določenim tipom, podobno kot razredi. V zgornjem primeru programa je lahko vse obnašanje, povezano s tipom `Student`, združeno z ustvarjanjem metod s prejemnikom tipa `Student`.

- Metode z istim imenom so dovoljene, če imajo različne tipe prejemnikov, medtem ko dveh funkcij z istim imenom v paketu ne moremo definirati. Recimo, da bi imeli v zgornji kodi dodatno strukturo `Profesor`, ki bi vsebovala polji `name` in `surname`. Za to strukturo lahko ustvarimo metodo `displayNameSurname`, ki že obstaja. Funkcije z istim imenom pa ne bi mogli.

Prejemnik metode je lahko vrednost ali kazalec. Če uporabimo vrednost, metoda dobi **kopijo** vrednosti, nad katero je bila poklicana, zato spremembe, ki jih naredi, niso vidne zunaj metode. Primer je metoda `setYearCopy(year int)`, ki spremeni le svojo kopijo študenta, zato letnik po njenem klicu ostane nespremenjen. Če pa kot prejemnika uporabimo kazalec na strukturo, metoda dostopa do izvirne vrednosti in spremembe se ohranijo tudi po zaključku metode. Primer take metode je `setYear(year int)`. Metoda `getYear()` vrednosti ne spreminja, zato je zanjo prejemnik tipa vrednost primeren.

## Vmesniki

V jeziku Go je vmesnik (angl. *interface*) nabor podpisov metod. Ko za nek podatkovni tip pripravimo definicijo vseh metod v vmesniku, pravimo, da podatkovni tip implementira vmesnik. To je zelo podobno svetu objektno usmerjenega programiranja. Vmesnik določa, katere metode mora tip imeti, in tip določi, kako so metode implementirane. Preko vmesnikov Go podpira polimorfizem v kodi.

Sintaksa vmesnika je naslednja:
```Go
type interfaceName interface {
        method1() type
        method2() type
        ...
}
```
Primer uporabe vmesnikov v jeziku Go:
```Go
/*
Program prikazuje uporabo vmesnikov v programskem jeziku Go
*/
package main

import (
	"errors"
	"fmt"
	"math"
)

// Definiramo nov vmesnik za izračun površine lika
// Če površine ni mogoče izračunati, metoda vrne napako
type areaCalculator interface {
	area() (float32, error)
}

// Definiramo strukturo za krog
type circle struct {
	radius float32
}

// Definiramo strukturo za pravokotnik
type rect struct {
	width  float32
	height float32
}

// Metoda za izračun površine kroga
func (c circle) area() (float32, error) {
	if c.radius < 0 {
		return 0, errors.New("polmer kroga ne sme biti negativen")
	}
	return c.radius * c.radius * math.Pi, nil
}

// Metoda za izračun površine pravokotnika
func (r rect) area() (float32, error) {
	if r.width < 0 || r.height < 0 {
		return 0, errors.New("stranici pravokotnika ne smeta biti negativni")
	}
	return r.width * r.height, nil
}

func main() {

	// Ustvarimo nekaj likov
	circle1, circle2, circle3 := circle{10}, circle{3}, circle{8}
	rect1, rect2, rect3 := rect{2, 2}, rect{3, 8}, rect{10, 10}
	// Ustvarimo še lik z neveljavnimi podatki
	invalid := circle{-1}

	// Ustvarimo rezino vmesnikov in vanjo damo vse like, ki implementirajo vmesnik
	shapes := []areaCalculator{circle1, circle2, circle3, rect1, rect2, rect3, invalid}

	// Izračunajmo skupno površino vseh likov v rezini
	totalArea := float32(0.0)
	for _, v := range shapes {
		a, err := v.area()
		// Like, katerih površine ni mogoče izračunati, preskočimo
		if err != nil {
			fmt.Println("Napaka:", err)
			continue
		}
		totalArea = totalArea + a
	}
	fmt.Printf("Skupna površina likov: %f\n", totalArea)
}
```
[vmesniki.go](koda/02-vmesniki/vmesniki.go)

Rezino `shapes` smo napolnili z različnimi liki. Ker vsi liki implementirajo vmesnik `areaCalculator`, lahko zelo enostavno izračunamo skupno površino vseh likov, ne glede na tip lika. Če bi v prihodnosti v kodo dodali še kak tip lika, katerega površino želimo računati, to ne bi zahtevalo večjih sprememb v kodi.

### Vmesnik error

Vmesniki so v jeziku Go prisotni tudi v samem jeziku in standardni knjižnici. Eden najpogostejših je vgrajeni vmesnik `error`, ki je običajen način za sporočanje napak:
```Go
type error interface {
	Error() string
}
```
Funkcija ali metoda, pri kateri lahko pride do napake, napako vrne kot zadnjo povratno vrednost. V zgornjem primeru metoda `area` poleg površine vrne še vrednost tipa `error`. Če so podatki lika neveljavni, ustvari napako s funkcijo `errors.New`, sicer pa namesto napake vrne `nil`. Klicatelj takoj po klicu preveri, ali je napaka različna od `nil`, in se odloči, kako jo bo obravnaval. V našem primeru napako izpiše in lik preskoči. Funkcije paketa `fmt` pri izpisu napake samodejno pokličejo njeno metodo `Error`.

### Prazen vmesnik

Vmesnik, ki nima nobene metode, se imenuje prazen vmesnik. Predstavljen je z `interface{}`. Ker prazen vmesnik nima nobenih metod, vsi podatkovni tipi implementirajo prazen vmesnik. Ker se prazni vmesniki v jeziku Go pogosto uporabljajo, imamo zanje definirano sopomenko `any`.

Primer:
```Go
/*
Program prikazuje uporabo praznih vmesnikov v programskem jeziku Go
*/
package main

import (
	"fmt"
)

// Funkcija kot argument prejme prazen vmesnik
// Namesto interface{} lahko pišemo tudi any
func display(x interface{}) {
	fmt.Printf("Tip: %T, Vrednost: %v\n", x, x)
}

// Funkcija z izbiro tipa različne podatkovne tipe obravnava različno
func describe(x any) {
	switch v := x.(type) {
	case string:
		fmt.Println("Niz dolžine", len(v))
	case int:
		fmt.Println("Celo število", v)
	default:
		fmt.Printf("Neznan tip: %T\n", v)
	}
}

func main() {
	// Uporaba funkcije display nad različnimi podatkovnimi tipi
	s := "PS"
	display(s)
	i := 42
	display(i)
	rect := struct {
		width  float32
		height float32
	}{
		width:  5,
		height: 4,
	}
	display(rect)

	// Preverjanje tipa vrednosti, shranjene v praznem vmesniku
	var x any = s
	str, ok := x.(string)
	fmt.Println(str, ok) // PS true
	n, ok := x.(int)
	fmt.Println(n, ok) // 0 false

	// Brez druge povratne vrednosti bi naslednji stavek povzročil paniko
	// m := x.(int)

	// Uporaba funkcije describe nad različnimi podatkovnimi tipi
	describe(s)
	describe(i)
	describe(rect)
}
```
[prazen-vmesnik.go](koda/03-prazen-vmesnik/prazen-vmesnik.go)

Iz vrednosti, shranjene v praznem vmesniku, izvirne vrednosti ne dobimo samodejno, saj prevajalnik njenega tipa ne pozna. Do nje pridemo s preverjanjem tipa (angl. *type assertion*), ki vrne vrednost in logično spremenljivko, ki pove, ali je vmesnik res vseboval zahtevani tip. Če tip ni pravi, dobimo ničelno vrednost zahtevanega tipa in vrednost `false`. Če drugo povratno vrednost izpustimo in je tip napačen, se program konča s paniko.
```Go
    var x any = s
    str, ok := x.(string)
    fmt.Println(str, ok) // PS true
    n, ok := x.(int)
    fmt.Println(n, ok) // 0 false
```
Kadar želimo obravnavati več možnih tipov, uporabimo izbiro tipa (angl. *type switch*), kot v funkciji `describe` v zgornjem primeru. Spremenljivka `v` ima v vsaki veji tip, naveden v tej veji, zato lahko v veji `string` kličemo `len(v)`, česar nad vrednostjo tipa `any` ne bi mogli narediti. V veji `default` ima `v` tip izraza `x`, torej prazen vmesnik.

# Domača naloga 3

Napišite program v jeziku Go, ki bo omogočal procesiranje naročil, ki jih kupci oddajo preko spletne trgovine. Rešitev oddajte preko [spletne učilnice](https://ucilnica.fri.uni-lj.si/mod/assign/view.php?id=54253).

**Navodila:**

Definirajte vmesnik `narocilo`:
```Go
type narocilo interface {
	obdelaj() error
}
```

in tri strukture za različne vrste naročil, ki jih obiskovalci spletne trgovine lahko izvedejo:

```Go
type izdelek struct {
	imeIzdelka string
	cena       float64
	teza       float64
}

type eknjiga struct {
	naslovKnjige string
	cena         float64
}

type spletniTecaj struct {
	imeTecaja   string
	trajanjeUre int
	cenaUre     float64
}
```

Za vsak tip naročila implementirajte vmesnik `narocilo`. Metoda `obdelaj` naj izpiše postavke naročila na zaslon in vrednost naročila prišteje globalni spremenljivki `promet`, v kateri beležimo skupno vrednost vseh obdelanih naročil. Vrednost naročila je pri izdelku in e-knjigi enaka ceni, pri spletnem tečaju pa zmnožku `trajanjeUre * cenaUre`. Prav tako naj se v globalni spremenljivki `stNarocil` beleži skupno število obdelanih naročil.

Metoda `obdelaj` naj najprej preveri, ali je naročilo veljavno (npr. ime ni prazno, cena, teža in trajanje so pozitivni). Če naročilo ni veljavno, naj vrne napako, ki pove, katero naročilo je neveljavno in zakaj, vrednosti spremenljivk `promet` in `stNarocil` pa naj ostanejo nespremenjene. Za sestavljanje sporočila o napaki lahko uporabite funkcijo `fmt.Errorf`, ki deluje podobno kot `fmt.Sprintf`.

Primer izpisa metode `obdelaj()`:
```
Številka naročila: 1
Ime izdelka: Prenosni računalnik
Cena: 2000 €
Teža: 2.5 kg
```

Ustvarite rezino in jo napolnite z naključnimi naročili različnih tipov z ustrezno inicializiranimi postavkami. V rezino dodajte vsaj 10 naročil, od tega vsaj dve neveljavni.

Za procesiranje vsakega naročila v rezini ustvarite svojo gorutino, ki naj naročilo obdela s klicem metode `obdelaj`. Stavek `go` povratne vrednosti zagnane funkcije zavrže, zato gorutina glavni gorutini napake ne more vrniti neposredno. Gorutine naj zato napake pošiljajo v skupen kanal. Ko se vse gorutine zaključijo, naj glavna gorutina zapre kanal in izpiše vse prejete napake.

Poskrbite, da glavna gorutina počaka na vse gorutine, ki obdelujejo naročila, nato naj izpiše število in skupni znesek obdelanih naročil. Pri posodabljanju spremenljivk `promet` in `stNarocil` s ključavnicami poskrbite, da ne pride do hkratnega spreminjanja vrednosti s strani več gorutin. Prav tako poskrbite, da se izpisi različnih naročil med sabo ne premešajo.

**Rok za oddajo: 9. 11. 2026**
