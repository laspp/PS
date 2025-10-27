# Uvod v programski jezik go

# Vzpostavitev okolja
Navodila za namestitev programskega jezika **go** najdete na [tukaj](https://go.dev/doc/install). Trenutna različica je 1.25.3. Programski jezik **go** je v učilnicah že nameščen. Če uporabljate [VSCode](https://code.visualstudio.com/) za programiranje, si lahko za lažje programiranje namestite [razširitev za go](https://code.visualstudio.com/docs/languages/go), ki nudi avtomatsko dopolnjevanje kode, dokumentacijo v živo in razhroščevalnik. Pri predmetu bomo uporabljali gručo [Arnes](https://www.arnes.si/arnes-z-najzmogljivejsim-superacunalnikom-v-sloveniji/), kjer je **go** nameščen v obliki modula, ki ga moramo naložiti, če ga želimo uporabljati. Na gruči je nameščena različica 1.23.6.

Na gruči Arnes naložimo modul `Go`, potem, ko smo se nanjo povezali preko protokola SSH, s sledečim ukazom:

```Bash
$ module load Go
```

Preverimo, če je **go** sedaj na voljo:

```Bash
$ go version
go version go1.23.6 linux/amd64
```


# Prvi program

Odpremo novo datoteko `gopher.go` in vanjo dodamo naslednji [izsek kode]((koda/gopher.go)):
```Go
/*
Program gopher izriše podobo ASCII ravninske mošnjičarke (angl. Pocket Gopher),
maskote programskega jezika go.
*/
package main

import "fmt"

func main() {
	//Prikaži ASCII Art
	fmt.Println("\nRavninska mošnjičarka (angl. Pocket gopher)\n")
	fmt.Println("       `.-::::::-.`")
	fmt.Println("   .:-::::::::::::::-:.")
	fmt.Println("   `_:::    ::    :::_`")
	fmt.Println("    .:( ^   :: ^   ):.")
	fmt.Println("    `:::   (..)   :::.")
	fmt.Println("    `:::::::UU:::::::`")
	fmt.Println("    .::::::::::::::::.")
	fmt.Println("    O::::::::::::::::O")
	fmt.Println("    -::::::::::::::::-")
	fmt.Println("    `::::::::::::::::`")
	fmt.Println("     .::::::::::::::.")
	fmt.Println("       oO:::::::Oo")

}
```
[gopher.go](koda/gopher.go)

Prav tako je priporočljivo ustvariti nov modul s katerim definiramo morebitne odvisnosti našega programa od drugih (nevgrajenih) paketov. To storimo tako, da v direktoriju, kjer se nahaja naš program, zaženemo ukaz:
```Bash
$ go mod init gopher
```
**Go** bo samodejno ustvaril datoteko `go.mod` z vsebino
```
module gopher

go 1.23.6
```
**Go** nudi več možnosti pri zaganjanju programa. Ker je prevajalnik relativno hiter, lahko program prevedemo in zaženemo z enim ukazom, to storimo na naslednji način:
```Bash
$ srun --reservation=fri go run gopher.go
```
Program smo prevedli in zagnali na enem od računskih vozlišč znotraj rezervacije `fri`. Opazimo lahko, da ukaz za sabo ni pustil nobene izvršljive datoteke, saj se je ta ustvarila samo začasno. Če hočemo ustvariti izvršljivo datoteko, uporabimo naslednji ukaz:
```Bash
$ go build gopher.go
```

Sedaj lahko naš program poženemo kot vsak drug izvršljiv program:
```Bash
$ srun --reservation=fri gopher
```

Uspešno smo zagnali naš prvi program napisan v **go**. Poglejmo si sedaj njegovo zgradbo.
Na vrhu datoteke je komentar, ki opisuje vsebino datokete. **Go** pozna dve vrsti komentarjev: vrstične, ki jih napovemo s sekvenco znakov `//` in bločne, ki se nahajajo med kombinacijo znakov `/* */`. 

Sledi definicija paketa, ki mu pripada datoteka z izvorno kodo `package main`. Datoteka, ki služi kot vstopna točka programa, to je, vsebuje funkcijo `main`, mora biti del paketa `main`. S pomočjo ključne besede `package` definiramo, kateremu paketu oziroma knjižnici pripada datoteka.

Sledi del kode, kjer napovemo uvoz morebitnih modulov oziroma paketov, ki jih bomo uporabljali v kodi. V danem primeru uvozimo paket `fmt`, ki vsebuje funkcije izpisovanje na standardni izhod, ki jih bomo kasneje potrebovali. Uvoženemu paketu lahko pripišemo tudi sopomenko, da si olajšamo pisanje kode. Primer:
```Go
import f "fmt"
```
Na ta način lahko v kodi namesto `fmt.Println()` napišemo kar `f.Println()`. Ko uvažamo več paketov, se lahko izognemo večkratnemu pisanju ključne besede `import` s pomočjo naslednje sintakse:
```Go
import(
    "fmt"
    "time"
)
```

Dokumentacijo o posameznem paketu ali funkciji, ki je del paketa, izpišemo z ukazom `go doc`. Primer:
```Bash
$ go doc fmt.Println
``` 

V kodi nato, podobno kot pri ostalih programskih jezikih, sledi definicija edine funkcije v datoteki `func main()`, ki služi kot vstopna točka našega programa. V njenem telesu imamo niz klicev funkcije `Println`, ki je del paketa `fmt`. Ena od posebnosti programskega jezika **go** je, kako v knjižnici oziroma paketu označimo, katere spremenljivke, funkcije, itd. so javne &mdash; jih lahko kličemo iz druge datoteke. To naredimo tako, da jih poimenujemo z veliko začetnico, kot v primeru `Println`. Imena napisana z malo začetnico pa ostanejo skrita (privatna). Pri poimenovanju spremenljivk funkcij, itd., se tipično uporablja format "camelCase".

V nadaljevanju si bomo pogledali nekaj osnovnih struktur programskega jezika **go**. Bolj podrobno dokumentacijo pa najdete na [spletni strani](https://go.dev/doc/) jezika. 

# Podatkovni tipi in spremenljivke

**Go** pozna naslednje podatkovne tipe:
```Go
bool //boolean

//predznačena in nepredznačena cela števila
int, int8, int16, int32, int64 // velikost int je odvisna od sistema
uint, uint8, uint16, uint32, uint64, uintptr

byte // sopomenka za uint8

rune // sopomenka za int32, uporablja se za predstavitev znakov unicode

// števila v plavajoči vejici
float32, float64

// kompleksna števila
complex64, complex128

// nizi
string
```

Primeri uporabe različnih podatkovnih tipov:

```Go
/*
Program tipi_spremenljivke prikazuje načine definiranja različnih osnovnih tipov spremenljivk v go.
*/
package main

import "fmt"

var b bool = true
var x int = 42
var y float64 = 3.141592653
var c complex128 = 42 + 3.14159i
var ch rune = 'Č'
var s string = "PS"

func main() {
	// v isti vrstici lahko definiramo več spremenljivk istega tipa
	var i, j int = 1, 2

	// operator := omogoča hkratno definiranje spremenljivke brez ključne besede var
	// podatkovni tip se določi samodejno glede na vrednost, ki jo prirejamo
	// dovoljeno samo znotraj funkcij
	k := 3
	banana, jabolko, brokoli := true, false, "Ne!"

	//Println samodejno uporabi ustrezno formatiranje glede na podatkovni tip
	fmt.Println(b, x, y, c, ch, s)
	fmt.Printf("%c\n", ch)

	fmt.Println(i, j, k, c, banana, jabolko, brokoli)
	// %T izpiše podatkovni tip posamezne spremenljivke
	fmt.Printf("%T %T %T %T %T %T %T %T\n", ch, i, j, k, c, banana, jabolko, brokoli)

	// Med podatkovnimi tipi lahko tudi pretvarjamo
	// Za razliko od programskega jezika C, je v go vedno potrebna eksplicitna konverzija pri priredbah med različnimi podatkovnimi tipi
	fx := float64(x)
	fmt.Printf("%T, %.1f\n", fx, fx)
}
```
[tipi_spremenljivke.go](koda/tipi_spremenljivke.go)

Posebnost programskega jezika **go** je operator `:=`, ki omogoča hitro definicijo nove spremenljivke, brez navedbe ključne besede `var` in podatkovnega tipa. Slednji se izpelje iz prirejene vrednosti. Ta operator lahko uporabimo samo znotraj funkcij. **Go** ima vgrajeno podporo za znake unicode znotraj nizov. Posamezni znak se obravnava kot podatkovni tip `rune`. V zgornji kodi smo uporabili funkcijo `Printf`, ki nudi formatirano izpisovaje na zaslon. Seznam podprtih formatov najdete na [povezavi](https://pkg.go.dev/fmt).

# Funkcije

Funkcije v **go** lahko vzamejo nič, enega ali več argumentov in vrnejo nič, enega ali več rezultatov; vrednosti, ki jih funkcija vrača lahko tudi poimenujemo. Jezik go podpira anonimne funkcije in zaprtja (*angl.* closures).

```Go
/*
Program funckije demonstrira uporabo funkcij v go.
Prikazani so različni načini podajanja argumentov, vračanja vrednosti, anonimne funkcije in zaprtja (closures)
*/
package main

import "fmt"

func sestej(a int, b int) int {
	return a + b
}

// Če ima več argumentov isti podatkovni tip, ga lahko pišemo samo enkrat, na koncu seznama argumentov
// Funkcija lahko vrača več vrednosti, podprto je tudi poimenovanje vračanih vrednosti
func deli(a, b int) (q, r int) {
	q = a / b
	r = a % b
	return
}

// Zaprtja funkcij (closures) omogočajo, da lahko znotraj funkcije naslavljamo spremenljivko, ki je definirana izven funkcije.
// Vrednost take spremenljive se ohranja tudi potem ko se iz funkcije vrnemo.
// Funkcija povecaj vrača anonimno funkcijo (zaprtje), ki spreminja spremenljivko i
func povecaj() func() {
	var i int
	return func() {
		i++
		fmt.Println(i)
	}
}

func main() {
	fmt.Println(sestej(42, 7))
	fmt.Println(deli(42, 7))
	// Go podpira anonimne funkcije, torej take, ki nimajo imena
	// Če anonimni funkciji sledi (), se le-ta takoj izvede
	result := func() int {

		fmt.Println("Anonimna funkcija")
		return 42
	}()
	fmt.Println(result)
	// Če pa ne, se spremenljivki priredi funkcija, ki jo lahko kasneje pokličemo
	f := func() {

		fmt.Println("Prirejena anonimna funkcija")
	}
	f()
	// Klic funckije, ki vrača zaprtje, vrednost i se skozi klice ohranja.
	p := povecaj()
	p()
	p()
}
```
[funkcije.go](koda/funkcije.go)
# Strukture
Strukture v **go** definiramo s pomočjo ključnih besed `type` in `struct`:
```Go
type imeStrukture struct {
  polje1 podatkovniTip
  polje2 podatkovniTip
  polje3 podatkovniTip
  ...
}
```
Primer uporabe struktur v programu in različni načini inicializacije:
```Go
/*
Program strukture prikazuje kako definiramo, inicializiramo in dostopamo do struktur v programskem jeziku go
*/
package main

import "fmt"

// Definicija strukture
type circle struct {
	x      int
	y      int
	radius int
	color  string
}

func main() {
	// Strukturo lahko inicializiramo na različne načine
	var smallCircle circle
	smallCircle.x = 0
	smallCircle.y = 0
	smallCircle.radius = 5
	smallCircle.color = "zelena"
	fmt.Println(smallCircle.x, smallCircle.y, smallCircle.radius, smallCircle.color)
	
	// Posamezna polja strukture lahko inicializiramo neposredno s pomočjo notacije {}
	bigCircle := circle{100, 100, 50, "rdeča"}
	fmt.Println(bigCircle)

	// Pri inicializaciji lahko navedemo imena polj, ki jim nastavljamo vrednost
	// ostala polja dobijo privzeto ničelno vrednost za dan podatkovni tip
	var mediumCircle = circle{radius: 15, color: "modra"}
	fmt.Printf("%T", mediumCircle)
}
```
[strukture.go](koda/strukture.go)
# Tabele in rezine
V programskem jeziku **go** poznamo dva tipa tabel: statične tabele (arrays) in rezine (slices), ki jih uporabljamo za ustvarjanje tabel spremenljive dolžine. Spremenljivko tipa tabela ustvarimo na naslednji način:
```Go
var T [n]podatkovniTip
```
S tem ustvarimo novo tabelo `T`, ki vsebuje `n` elementov tipa `podatkovniTip`. Dolžine tabele ne moremo spreminjati med izvajanjem.
```Go
/*
Program tabele prikazuje kako definiramo, inicializiramo in dostopamo do tabel v programskem jeziku go
*/
package main

import "fmt"

func main() {
	// Definicija in izpis tabele z dvema elementoma
	// Tabela ima vedno definirano velikost
	var a [2]string
	a[0] = "Programiram"
	a[1] = "Go"
	fmt.Printf("%T, %s, %s\n", a, a[0], a[1])
	fmt.Println(a)

	// Uporabimo lahko tudi kratko notacijo
	fibonacci := [6]int{1, 1, 2, 3, 5, 8}
	fmt.Println(fibonacci)

	// Večdimenzionalne tabele
	magic := [3][3]int{{2, 7, 6}, {9, 5, 1}, {4, 3, 8}}
	fmt.Println(magic)
}
```
[tabele.go](koda/tabele.go)

Uporaba rezin je v **go**-ju bolj pogosta kot uporaba tabel. Rezino definiramo podobno kot tabelo, le da ne navedemo velikosti.
```Go
var []podatkovniTip
```
Neinicializirana rezina je dolga 0 elementov in je enaka vrednosti `nil`.

Rezino lahko inicializiramo na več načinov. Lahko jo uporabimo kot referenco v obstoječo statično tabelo, ali pa jo ustvarimo s pomočjo vgrajene funkcije `make`.
```Go
/*
Program rezine prikazuje kako definiramo, inicializiramo in dostopamo do rezin v programskem jeziku go
*/
package main

import "fmt"

func main() {

	// fibonaci je tabela (array)
	fibonacci := [6]int{1, 1, 2, 3, 5, 8}
	// s1 je rezina (slice), ki služi kot referenca na del tabele fibonacii
	// z notacijo [:] povemo, kateri del tabele sestavlja rezino
	// nismo ustvarili kopije podatkov ampak samo referenco na elemente tabele
	var s1 []int = fibonacci[0:3]
	// Do elementov v rezini dostopamo na enak način kot do elementov tabele
	fmt.Printf("%T, %d, %d, %d\n", s1, s1[0], s1[1], s1[2])

	// Zgornjo ali spoodnjo ali obe meji lahko izpustimo, če ustvarjamo rezino, ki vsebuje vse elemente tabele
	var s2 []int = fibonacci[:]
	s2[0] = 0
	fmt.Println(s2)
	fmt.Println(fibonacci)

	// Rezino lahko ustvarimo tudi neposredno
	// hkrati se ustvari tabela s podatki in rezina, ki kaže na tabelo
	letters := []string{"a", "b", "c", "d"}
	fmt.Println(letters)

	// Vsaka rezina ima definirano dolžino in kapaciteto
	fmt.Printf("Dolzina s1:%d, Kapaciteta s1:%d\nDolzina s2:%d, Kapaciteta s2:%d\n", len(s1), cap(s1), len(s2), cap(s2))

	// Dolžino rezine lahko povečamo, če je na voljo dovolj kapacitete
	// Pri tem ne pride do odvečnega kopiranja podatkov
	s1 = s1[:4]
	fmt.Printf("Dolzina s1:%d | %d, %d, %d %d\n", len(s1), s1[0], s1[1], s1[2], s1[3])

	// Rezine lahko dinamično ustvarjamo s pomočjo funkcije make
	// Rezina s3 bo inicializirana z ničelnimi vrednosti in bo imela dolžino ter kapaciteto enako 5
	s3 := make([]int, 5)
	fmt.Println(s3)

	// Dodatno lahko specificiramo tudi kapaciteto rezine
	s4 := make([]int, 5, 10)
	fmt.Printf("Dolzina s4:%d, Kapaciteta s4:%d\n", len(s4), cap(s4))

	// Poznamo tudi večdimenzionalne rezine
	magic := [][]int{
		[]int{2, 7, 6},
		[]int{9, 5, 1},
		[]int{4, 3, 8},
	}
	fmt.Println(magic)

	// Rezini lahko dodajamo elemete s funkcijo append
	// Dolžina in kapaciteta rezine se povečata
	fmt.Printf("Dolzina letters:%d, Kapaciteta letters:%d | %v\n", len(letters), cap(letters), letters)
	letters = append(letters, "e", "f")
	fmt.Printf("Dolzina letters:%d, Kapaciteta letters:%d | %v\n", len(letters), cap(letters), letters)

	// Rezine lahko kopiramo
	copy(s3, s1)
	fmt.Printf("Dolzina s3:%d, Kapaciteta s3:%d |%v -> %v\n", len(s3), cap(s3), s1, s3)
}
```
[rezine.go](koda/rezine.go)
# Slovar
**Go** ima vgrajeno podatkovno strukturo slovar (map), ki deluje po principu zgoščene tabele. Do vrednosti v slovarju lahko dostopamo preko ključa. Nov slovar ustvarimo s pomočjo vgrajene funkcije `make` ali pa neposredno, kot bomo videli na primeru.
```Go
var slovar = make(map[podatkovniTipKljuca]PodatkovniTipVrednosti)
```
Primer uporabe slovarja:
```Go
/*
Program slovar prikazuje kako definiramo, inicializiramo in dostopamo do podatkovne strukture slovar (map) v programskem jeziku go
*/
package main

import "fmt"

func main() {
	// Definiramo in inicializiramo nov slovar
	var fruitSupply map[string]int
	// Dokler slovarja ne inicializiramo z make, do njega ne moremo dostopati
	fruitSupply = make(map[string]int)

	// Dodamo nove ključe v slovar
	fruitSupply["jabolka"] = 5
	fruitSupply["hruske"] = 3
	fmt.Println(fruitSupply)
	// Za ključe, ki še niso v slovarju, dobimo ničelno vrednost
	fmt.Println(fruitSupply["jabolka"], fruitSupply["slive"])

	// Izbrišemo ključ
	delete(fruitSupply, "jabolka")

	//Preverimo, če je ključ v slovarju
	e, ok := fruitSupply["hruske"]
	fmt.Println(e, ok)

	// Pobrišemo vse ključe
	clear(fruitSupply)
	fmt.Println(fruitSupply)

	// Ustvarimo nov slovar z zalogo sadja z neposredno inicializacijo ključev
	fruitSupply = map[string]int{"jabolka": 1, "hruske": 2, "slive": 10}
	fmt.Println(fruitSupply)
}
```
[slovar.go](koda/slovar.go)
# Zanke
**Go** pozna samo eno vrsto zanke in to je zanka `for`. Sestavljen je iz treh komponent: inicializacije, ki se zgodi pred prvo iteracijo; pogoja, ki se preveri ob vsaki iteraciji; in stavka, ki se izvede po vsaki iteraciji. Za razliko od programskega jezika C in podobnih, navedene tri komponente niso obdane z oklepaji. Vse tri komponente so opcijske. Na ta način dobimo različne tipe zank. Dodatno lahko poenostavimo pisanje zank z uporabo ključne besede `range`.

``` Go
for inicializacija; pogoj; stavek {
	// Telo zanke
}
```
Primeri zank for:
```Go
/*
Program zanke prikazuje kako uporabljamo zanko for v jeziku go
*/
package main

import "fmt"

func main() {
	// Tipičen primer zanke
	for i := 0; i <= 9; i++ {
		fmt.Println(i)
	}

	// Izpustimo inicializacijo in stavek po koncu iteracije (stavek while)
	// Podpičij nam ni treba pisati
	i := 0
	for i <= 9 {
		fmt.Println(i)
		i++
	}

	// Izpustimo še pogoj in dobimo neskončno zanko
	for {
		fmt.Println("Zanka")
		break
	}

	// Bolj elegantna oblika zanke for z uporabo range
	// Uporabno za sprehajanje po tabelah, rezinah, nizih, slovarjih
	fibonacci := [6]int{1, 1, 2, 3, 5, 8}
	for i, v := range fibonacci {
		fmt.Println(i, v)
	}
	fmt.Println()
	// Indeks ali vrednost lahko izpustimo
	for _, v := range fibonacci {
		fmt.Println(v)
	}
	fmt.Println()
	for i := range fibonacci {
		fmt.Println(i)
	}

	// Pri nizih obstaja razlika
	// Če se po nizu sprehodimo z range nam vrača rune (znaki unicode)
	// Če se po nizu sprehodimo z indeksom nam vrača posamezne bajte
	s := "Čuri Muri, kljukec Juri"
	for _, c := range s {
		fmt.Printf("%c", c)
	}
	fmt.Println()
	for i := 0; i < len(s); i++ {
		fmt.Printf("%c", s[i])
	}
	fmt.Println()
}
```
[zanke.go](koda/zanke.go)
# Izbirni stavki
Izbirna stavka `if-else` in `switch` delujeta podobno kot smo navajeni iz drugih programskih jezikov. Stavek  `if-else` ima naslednjo sintakso:
```Go
if stavek; pogoj1 {
	 // izpolnjen pogoj1
} else if pogoj2 {
	// izpolnjen pogoj2
} else {
	//sicer
}
```
Kombiniramo lahko več zaporednih `if-else`. Poleg pogoja je možno v `if` dodati tudi dodaten stavek, ki ponavadi služi inicializaciji spremenljivk.
Primeri `if-else` stavkov:
```Go
/*
Program if-else prikazuje kako uporabljamo pogoj if v jeziku go
*/
package main

import "fmt"

func main() {
	// Tipičen primer stavka if-else
	// Zaviti oklepaji so obvezni
	if 3%2 == 0 {
		fmt.Println("3 je sodo število")
	} else {
		fmt.Println("3 je liho število")
	}

	// Stavek if-else if-else z inicializacijo
	if num := 42; num < 0 {
		fmt.Println(num, "je negativno število")
	} else if num > 0 {
		fmt.Println(num, "je pozitivno število")
	} else {
		fmt.Println(num, "je enako nič")
	}
}
```
[if-else.go](koda/if-else.go)

Stavek `switch` lahko uporabimo, kadar bi imeli z uporabo `if-else` veliko število vejitev in s tem nepregledno kodo. Sintaksa `switch` je naslednja:
```Go
switch stavek; izraz {
		case izraz1: stavek...
		case izraz2: stavek...
		...
		default: stavek 
}
```
Začetni stavek, ki običajno služi inicializaciji, in izraz lahko tudi izpustimo. Možnost `default` se izvede, če noben drug izraz ni izpolnjen.
Primeri uporabe stavka `switch`:
```Go
/*
Program switch prikazuje kako uporabljamo stavek switch v jeziku go
*/
package main

import (
	"fmt"
	"time"
)

func main() {

	// Tipičen primer stavka switch
	i := 2
	switch i {
	case 1:
		fmt.Println("ena")
	case 2:
		fmt.Println("dva")
	case 3:
		fmt.Println("tri")
	default:
		fmt.Println("nič od naštetega")
	}
	// Izpuščen izraz: ekvivalentno stavku if-else
	i = 3
	switch {
	case i == 1:
		fmt.Println("ena")
	case i == 2:
		fmt.Println("dva")
	case i == 3:
		fmt.Println("tri")
	default:
		fmt.Println("nič od naštetega")
	}
	// Primer inicializacije
	// Izraze lahko povežemo
	switch dayOfTheWeek := time.Now().Weekday(); dayOfTheWeek {
	case time.Saturday, time.Sunday:
		fmt.Println("Vikend je!")
	default:
		fmt.Println("Na FRI bo treba!")
	}
}
```
[switch.go](koda/switch.go)
# Zakasnjeno izvajanje
**Go** nudi možnost zakasnjenega izvajanja s pomočjo stavka `defer`. Z njim označimo stavek/funkcijo, ki naj se izvede, ko se trenutna funkcija zaključi. Morebitni argumenti tako zakasnjeni funkciji se izračunajo takoj, medtem ko se sam klic zgodi zakasnjeno. Stavek `defer` po navadi uporabimo, kadar želimo že vnaprej napovedati, da bomo kako podatkovno strukturo počistili, ko se funkcija zaključi (npr. zapremo datoteko). S tem se izognemo morebitnim hroščem v kodi, saj bo **Go** poskrbel, da se datoteka zapre, ne glede na to kaj se med samim izvajanjem funkcije dogaja.
```Go
/*
Program zakasnjeno-izvajanje prikazuje uporabo stavka defer v jeziku go
*/
package main

import "fmt"

func main() {
	// Vrednost s se bo določila takoj, izpis pa se bo zgodil na koncu
	s := "svet"
	defer fmt.Println(s)
	s = "Pozdravljen"
	fmt.Println(s)

	//Klici zakasnjenih funkcij se nalagajo na sklad in se izvedejo po principu LIFO
	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}
}

```   

# Kazalci
V **go** poznamo tudi kazalce. Uporablja se podobna sintaksa kot v programskem jeziku C.
```Go
var p *podatkovniTip
```
Operator * nam vrne vrednost na katero kaže kazalec, operator & pa vrne naslov spremenljivke. Neinicializiran kazalec dobi privzeto vrednost `nil`.
```Go
/*
Program kazalci prikazuje uporabo kazalcev v jeziku go
*/
package main

import "fmt"

type point struct {
	x int
	y int
}

func main() {
	i, j := 42, 1337
	// Ustvarimo kazalec p, ki kaže na spremenljivko i
	p := &i
	// Dostop do vrednosti, na katero kaže kazalec
	fmt.Println(*p)
	// Spremenimo vrednost preko kazalca
	*p = j
	fmt.Println(i)

	// Pri kazalcih na strukture lahko uporabimo poenostavljeno sintakso brez *, ko spreminjamo vrednosti preko kazalca
	t := point{}
	p2 := &t
	p2.x = 15
	fmt.Println(t)
}
```
[kazalci.go](koda/kazalci.go)

# Naprednejši konstrukti (metode, vmesniki, parametrizirani tipi, ...)
Naprednejše konstrukte kot so metode vmesniki, parametrizirani tipi bomo zaenkrat izpustili. Kogar zanimajo naprednejše funkcinalnosti, jih lahko razišče s pomočjo [spletne dokumentacije](https://go.dev/doc/) in mnogih vodičev. 

Priporočam, da se sprehodite skozi vodič [Tour of Go](https://go.dev/tour/welcome/1) in predelajte primere ter vaje vsaj do poglavja Methods.

# Domača naloga 1

Napišite program v jeziku Go, ki bo omogočal dodajanje ocen študentov za posamezne domače naloge in izračun končne ocene nalog. Rešitev oddajte preko [spletne učilnice](https://ucilnica.fri.uni-lj.si/mod/assign/view.php?id=37715).

**Navodila:**

Definirajte strukturo Student:
```Go
type Student struct {
    ime     string
    priimek string
    ocene   []int
}
```
in ustvarite **slovar** študentov (struktur `Student`). Ključ naj bo vpisna številka študenta. V slovar dodajte vsaj tri študente s poljubnimi ocenami.

Napišite naslednje funkcije:
- `dodajOceno(studenti map[string]Student, vpisnaStevilka string, ocena int)`, ki študentu z dano vpisno številko doda v seznam ocen novo oceno. Pred dodajanjem ocene naj funkcija preveri, ali je ocena v ustreznem območju (med 0 in 10). Če študent ne obstaja, naj funkcija izpiše obvestilo, da študenta ni na seznamu.
 - `povprecje(studenti map[string]Student, vpisnaStevilka string) float64` vrne povprečno oceno študenta. Če študent ne obstaja, naj funkcija vrne `-1.0`. Če je število ocen domačih nalog manjše od `6`, funkcija vrne povprečje `0.0`.
 - `izpisRedovalnice(studenti map[string]Student)`, ki izpiše imena in priimke vseh študentov in njihove ocene.
	Primer izpisa:
	```
	REDOVALNICA:
	63210001 - Ana Novak: [10 9 8]
	63210002 - Boris Kralj: [6 7 5 8]
	63210003 - Janez Novak: [4 5 3 5]
	```
 - `izpisiKoncniUspeh(studenti map[string]Student)`, ki za vsakega študenta izpiše njegovo ime in priimek ter povprečno oceno. Znotraj funkcije uporabite funkcijo `povprecje` Poleg tega naj za posameznega študenta izpiše tudi komentar:
  
   - če je povprečna ocena >=9 naj izpiše `Odličen študent!`;
   - če je povprečna ocena med 6 in 9, naj izpiše `Povprečen študent`;
   - če je povprečna ocena <6, naj izpiše `Neuspešen študent`.
 
	Primer izpisa:
	```
	Ana Novak: povprečna ocena 9.0 -> Odličen študent!
	Boris Kralj: povprečna ocena 6.0 -> Povprečen študent
	Janez Novak: povprečna ocena 4.0 -> Neuspešen študent
	```

V vašem programu demonstrirajte delovanje vseh štirih funkcij!

**Rok za oddajo: 2. 11. 2025**
