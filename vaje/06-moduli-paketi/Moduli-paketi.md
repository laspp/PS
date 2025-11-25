# Paketi in moduli v Go

Upravljanje s knjižnicami in kodo se v Go vrti okoli treh konceptov:  repozitorijev (*repositories*), modulov (*modules*) in paketov (*packages*). Uporabo repozitorijev kode kot sta [GitHub](https://github.com/) in [GitLab](https://gitlab.com/) dobro poznamo tudi iz drugih programskih jezikov. Namenjeni so hranjenu kode znotraj sistema za upravljanje z izvorno kodo. Modul v Go predstavlja projekt (aplikacijo ali knjižnico), ki ga hranimo v repozitoriju. Moduli so sestavljeni iz enega ali več paketov, ki so osnovna organizacijska enota kode v Go s pomočjo katerih smiselno strukturiramo kodo znotraj modula.  

## Paketi

Paket je osnovna organizacijska enota kode v Go. Sestavlja jo ena ali več datotek znotraj iste mape, ki skupaj rešujejo določen problem. Paket ustvarimo tako, da znotraj datoteke `.go` navedemo:

```go
package imepaketa
```

Imenu paketa sledi koda, ki pripada danemu paketu.
Primer:
```go
package mymath

func Square(a int) int{
    return multiply(a,a)
}

func mutliply(a int, b int) int{
    return a*b
}

```
Simbole (spremenljivke, funkcije, tipe, ...) znotraj paketa naredimo vidne ostalim delom knjižnice ali aplikacije, tako da jih pišemo z **Veliko začetnico**. V zgornjem primeru paketa `mymath` je funkcija `multiply` privatna (ni vidna in je ne moremo poklicati izven paketa `mymath`). Pakete hranimo v mapah, ki imajo isto ime kot paket. 
Paket `main` je poseben, saj ga mora vsebovati vsaka aplikacija. V njem se nahaja vstopna točka aplikacije (funkcija `main`). Knjižnice paketa `main` ne potrebujejo, saj ne vsebujejo vstopne točke.

Primeri standardnih paketov: `fmt`, `math`, `net/http`.
Paket uvozimo s pomočj stavka `import`:

```go
import (
    "fmt"
    "math"
    "mymath"
    http "net/http" // Uporaba sopomenke (alias), da skrajšamo ime paketa v kodi
    "github.com/davors/weatherstation/weather" // Uvoz paketa weather, ki se nahaja v repozitoriju
)

func main() {
    fmt.Println(mymath.Square(4))
    fmt.Println(mymath.multiply(4,4)) // Napaka, funckija multiply ni vidna izven paketa mymath
    //...
}
```
## Moduli

Module predstavlja aplikacijo ali knjižnico, ki jo razvijamo in ga običajno hranimo v repozitoriju.
Ustvarimo ga z ukazom:

```bash
go mod init ime_modula
```

Tipično je ime modula enako repozitoriju (GitHub/Gitlab) v katerem ga bomo hranili. Primer: 

```bash
go mod init github.com/davors/weatherstation
```

Modul je definiran z datoteko `go.mod`, ki vsebuje ime modula in odvisnosti.
Datoteka `go.mod` omogoča:
 - upravljanje odvisnosti (knjižnice in verzije),
 - reproducibilno gradnjo projekta,
 - gradnjo projektov, ki vsebujejo več paketov.

Struktura modula tipično izgleda tako:
```
ime_modula/        ← modul
  go.mod
  main.go          ← paket main
  paket1/      ← paket 1
    p1.go
  paket2/      ← paket 2
    p2.go
```
Če je struktura modula bolj kompleksna (veliko paketov in izvršljivih datotek) se v Go poslužujemo naslednjega pristopa pri organizaciji kode:

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
```
Izvršljivo datoteko bi iz kode v datoteki `exe1.go` ustvarili z ukazom:
```
go build ./cmd/exe1
```
Pri prevajanju aplikacije s pomočjo ukaza `go build` bo Go samodejno poiskal vse odvisnosti, in jih po potrebi prenesel in prevedel. Če v modulu, ki ga razvijamo uporabljamo pakete, ki se nahajajo v zunanjem repozitoriju, lahko nastavimo odvisnosti v datoteki `go.mod` s pomočjo ukazov:
```bash
go mod tidy
```
ali
```bash
go get ime_paketa
```
Googlov strežnik (proxy.golang.org) samodejno indeksira in hrani module Go, ki jih najde v javnih repozitorijih. Go, če ga ne drugače nastavimo, pakete uvaža preko strežnika proxy in ne neposredno iz repozitorija. Googlovi strežniki samodejno zbirajo dokumentacijo o javnih modulih in jo objavljajo na spletni strani https://pkg.go.dev/.

# Domača naloga 5

Napišite modul v jeziku Go, in ga objavite na enem od javnih repozitorijev.
Preko spletne [spletne učilnice](https://ucilnica.fri.uni-lj.si/mod/assign/view.php?id=60604) oddajte povezavo do modula.

**Navodila:**

Kodo iz [domače naloge 1](../02-programski-jezik-go/Uvod_v_go.md#domača-naloga-1) zapakirajte v modul in jo objavite v javnem repozitoriju.
Znotraj modula ustvarite paket `main`, ki bo uporabljal funkcije/metode paketa `redovalnica`.

Paket `redovalnica` naj izvaža naslednje funkcije/metode: 
- `DodajOceno`, 
- `IzpisVsehOcen`
- `IzpisiKoncniUspeh`

Funkcija `povprecje` naj bo del paketa vendar naj ostane skrita.

Kodo rešitve 1. domače naloge lahko poljubno reorganizirate. Le funkcionalnost izvoženih funkcij/metod naj ostane enaka. V vaši rešitvi uporabite paket "github.com/urfave/cli/v3" za gradnjo aplikacij z ukaznim vmesnikom. Aplikaciji dodajte tri stikala:
 - `stOcen`, ki definira najmanjše število ocen potrebnih za pozitivno oceno;
 - `minOcena`, najmanjša možna ocena;
 - `maxOcena`, največja možna ocena.

Zgledujete se lahko po [primeru](https://github.com/davors/weatherstation). Paketu `redovalnica` dodajte ustrezno dokumentacijo. Pri pisanju dokumentacije se poskušajte držati [pravil](https://go.dev/blog/godoc).

**Rok za oddajo: 23. 11. 2025**
