# Paketi in moduli v Go

Upravljanje s knjižnicami in kodo se v Go vrti okoli treh konceptov: repozitorijev (angl. *repositories*), modulov (angl. *modules*) in paketov (angl. *packages*). Uporabo repozitorijev kode, kot sta [GitHub](https://github.com/) in [GitLab](https://gitlab.com/) dobro poznamo tudi iz drugih programskih jezikov. Namenjeni so hranjenju kode znotraj sistema za upravljanje z izvorno kodo. Modul v Go predstavlja projekt (aplikacijo ali knjižnico), ki ga hranimo v repozitoriju. Moduli so sestavljeni iz enega ali več paketov, ki so osnovna organizacijska enota kode v Go in s pomočjo katerih smiselno strukturiramo kodo znotraj modula.

## Paketi

Paket je osnovna organizacijska enota kode v Go. Sestavlja ga ena ali več datotek znotraj iste mape, ki skupaj rešujejo določen problem. Paket ustvarimo tako, da znotraj datoteke `.go` navedemo:

```go
package imepaketa
```

Imenu paketa sledi koda, ki pripada danemu paketu.
Primer:
```go
package mymath

func Square(a int) int {
	return multiply(a, a)
}

func multiply(a int, b int) int {
	return a * b
}
```
Simbole (spremenljivke, funkcije, tipe, ...) znotraj paketa naredimo vidne drugim paketom, tako da jih pišemo z veliko začetnico. Pravimo, da jih paket izvozi. V zgornjem primeru paketa `mymath` je funkcija `multiply` neizvožena (ni vidna in je ne moremo poklicati izven paketa `mymath`). Pakete običajno hranimo v mapah, ki imajo isto ime kot paket.
Paket `main` je poseben, saj ga mora vsebovati vsaka aplikacija. V njem se nahaja vstopna točka aplikacije (funkcija `main`). Knjižnice paketa `main` ne potrebujejo, saj ne vsebujejo vstopne točke.

Primeri standardnih paketov: `fmt`, `math`, `net/http`. Paket uvozimo s pomočjo stavka `import`, ki mu podamo pot uvoza paketa. Pot uvoza paketa iz modula pa je sestavljena iz imena modula in poti do mape paketa znotraj modula. Če se paket `mymath` nahaja v mapi `mymath` modula z imenom `mymodule`, ga uvozimo s potjo `mymodule/mymath`, v kodi pa ga uporabljamo pod imenom, ki je naveden v stavku `package`:

```go
import (
	"fmt"
	"math"

	"github.com/davors/weatherstation/weather" // Uvoz paketa weather iz zunanjega modula
	"mymodule/mymath"                          // Uvoz paketa mymath iz našega modula mymodule
)

func main() {
	fmt.Println(mymath.Square(4))
	fmt.Println(mymath.multiply(4, 4)) // Napaka, funkcija multiply ni vidna izven paketa mymath
	//...
}
```
## Moduli

Modul predstavlja aplikacijo ali knjižnico, ki jo razvijamo in običajno hranimo v repozitoriju.
Ustvarimo ga z ukazom:

```bash
go mod init ime_modula
```

Običajno je ime modula enako poti do repozitorija (GitHub/GitLab), v katerem ga bomo hranili. Primer:

```bash
go mod init github.com/davors/weatherstation
```

Ime modula mora biti enako poti do repozitorija le, če želimo, da lahko modul z ukazom `go get` uvozijo tudi drugi. Za programe, ki jih ne objavljamo, je ime modula lahko poljubno. Primeri iz prejšnjih poglavij so bili tako urejeni v modul z imenom `koda`.

Modul je definiran z datoteko `go.mod`, ki vsebuje ime modula in odvisnosti.
Datoteka `go.mod` omogoča:
 - upravljanje odvisnosti (paketi in njihove različice),
 - ponovljivo gradnjo projekta,
 - gradnjo projektov, ki vsebujejo več paketov.

Struktura modula je običajno takšna:
```
ime_modula/        ← modul
  go.mod
  main.go          ← paket main
  paket1/      ← paket 1
    p1.go
  paket2/      ← paket 2
    p2.go
```
Če je struktura modula bolj kompleksna (veliko paketov in izvršljivih datotek), se v Go poslužujemo naslednjega pristopa pri organizaciji kode:

```
ime_modula/        ← modul
  go.mod
  cmd/
    exe1/
        exe1.go    ← paket main
    exe2/
        exe2.go    ← paket main
  pkg/
    paket1/        ← paket 1
        p1.go
    paket2/        ← paket 2
        p2.go
  internal/
    paket3/        ← paket 3, viden le znotraj modula
        p3.go
```
Izvršljivo datoteko bi iz kode v datoteki `exe1.go` ustvarili z ukazom:
```bash
go build ./cmd/exe1
```
Izvršljiva datoteka dobi ime po mapi, v kateri se nahaja paket `main`, v tem primeru `exe1`. Tako je urejen tudi [primer vremenske postaje](koda/), kjer se program nahaja v mapi `cmd/weatherstation`, paket `weather` pa neposredno v korenu modula.

Primer zaženemo iz mape `koda` z ukazom:
```bash
go run ./cmd/weatherstation
```
Aplikacija za branje stikal uporablja paket `github.com/urfave/cli/v3` in podpira dve stikali:
 - `--pollTime`, interval osveževanja senzorjev v milisekundah (privzeto 500);
 - `--timeout`, čas v sekundah, ki ga aplikacija čaka na meritev, preden javi napako (privzeto 5).

Stikala podamo za imenom programa, na primer `go run ./cmd/weatherstation --pollTime 200 --timeout 2`. Seznam vseh stikal izpiše stikalo `--help`. Aplikacijo zapremo s tipko `q`.

Mapa `internal` ima v Go poseben pomen. Paket v njej lahko uvozijo le paketi iz drevesa map, ki se začne v mapi nad mapo `internal`. Če se `internal` nahaja v korenu modula, ga lahko uvozijo vsi paketi istega modula, drugi moduli pa ne. V mapo `internal` zato postavimo pakete, ki jih potrebujemo znotraj modula, ne želimo pa, da jih uporabljajo drugi. 

Pri prevajanju aplikacije z ukazom `go build` Go po potrebi prenese in prevede le odvisnosti, ki so navedene v datoteki `go.mod`. Če v kodi uvozimo paket iz zunanjega modula, ki ga v `go.mod` še ni, se prevajanje ustavi z napako:
```
no required module provides package github.com/urfave/cli/v3; to add it:
	go get github.com/urfave/cli/v3
```
Odvisnost v datoteko `go.mod` dodamo z ukazom `go get`, ki mu podamo pot modula ali paketa in po želji še različico:
```bash
go get github.com/urfave/cli/v3@latest
```
ali z ukazom `go mod tidy`, ki pregleda uvoze v vseh paketih modula, doda manjkajoče odvisnosti in odstrani tiste, ki jih ne uporabljamo več:
```bash
go mod tidy
```

Če ne nastavimo drugače, Go module prenaša prek Googlovega strežnika proxy (proxy.golang.org) in ne neposredno iz repozitorija. Modul prenese iz repozitorija in shrani šele, ko ga nekdo prvič zahteva, na primer z ukazom `go get`. Dokumentacijo modulov, ki jih je strežnik proxy shranil, objavlja spletna stran [https://pkg.go.dev/](https://pkg.go.dev/).
Po zahtevi lahko traja nekaj minut, da se modul pojavi na strani pkg.go.dev. Dokumentacija se prikaže le, če repozitorij vsebuje datoteko z eno od prepoznanih odprtokodnih licenc (npr. MIT ali BSD).

## Dokumentacija modulov in paketov

Go ima vgrajen sistem za dokumentacijo, ki temelji na komentarjih v kodi. Komentar, ki stoji neposredno (brez prazne vrstice) pred stavkom `package` ali pred deklaracijo izvoženega imena, postane njegova dokumentacija. Po dogovoru je zapisan v celih povedih in se začne z imenom, ki ga opisuje. Primer iz paketa [weather](koda/weather/weather.go):
```go
// Package weather provides a simulated weather station that generates random weather data.
package weather

// NewStation creates a new weather station with the specified polling interval and timeout.
func NewStation(pollingInterval time.Duration, timeout time.Duration) *Station {
```
Dokumentacijo paketa pregledamo z ukazom `go doc`, ki mu podamo pot do paketa in po želji še ime, ki nas zanima:
```bash
go doc ./weather
go doc ./weather NewStation
```
Ista dokumentacija se za objavljene module prikaže tudi na strani pkg.go.dev. Pravila oblikovanja so opisana v [navodilih za pisanje dokumentacije](https://go.dev/doc/comment).

# Domača naloga 5

Napišite modul v jeziku Go in ga objavite na enem od javnih repozitorijev.
Preko [spletne učilnice](https://ucilnica.fri.uni-lj.si/mod/assign/view.php?id=60725) oddajte **povezavo do modula**.

**Navodila:**

Kodo iz [domače naloge 1](../02-programski-jezik-go/Uvod_v_go.md#domača-naloga-1) zapakirajte v modul in jo objavite v javnem repozitoriju.
Znotraj modula ustvarite paket `vreme`, ki vsebuje strukturo `Postaja` in funkcije za delo z vremenskimi postajami, ter paket `main`, ki te funkcije uporablja.

Paket `vreme` naj izvaža naslednje funkcije:
- `DodajMeritev`,
- `IzpisMeritev`,
- `IzpisPorocila`.

Funkcija `statistika` naj bo del paketa, vendar naj ostane skrita. Ker so polja strukture `Postaja` neizvožena, paket `main` postaj ne more ustvariti sam. Paketu zato dodajte še izvoženo funkcijo `DodajPostajo`, ki v slovar doda novo postajo z dano oznako, imenom in krajem.

Kodo rešitve 1. domače naloge lahko poljubno reorganizirate in po potrebi prilagodite parametre funkcij. Le funkcionalnost izvoženih funkcij naj ostane enaka. V vaši rešitvi uporabite paket ```github.com/urfave/cli/v3``` za gradnjo aplikacij z ukaznim vmesnikom. Aplikaciji dodajte tri stikala:
 - `minTemp`, najnižja dovoljena temperatura meritve (privzeto -50 °C);
 - `maxTemp`, najvišja dovoljena temperatura meritve (privzeto 50 °C);
 - `stMeritev`, najmanjše število meritev, potrebnih za izračun statistike (privzeto 3).

Zgledujete se lahko po [primeru](https://github.com/davors/weatherstation). Paketu `postaje` dodajte ustrezno dokumentacijo. Pri pisanju dokumentacije se poskušajte držati [pravil](https://go.dev/doc/comment).

Ko je modul pripravljen, v repozitoriju ustvarite oznako različice (npr. `v1.0.0`), tako da ga lahko drugi uvozijo z ukazom `go get <pot modula>@v1.0.0`. Preverite, ali se modul s to različico prikaže na strani [pkg.go.dev](https://pkg.go.dev/).

**Rok za oddajo: 23. 11. 2026**
