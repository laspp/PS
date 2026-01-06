# Programiranje grafičnih procesnih enot CUDA

Na gruči Arnes je na voljo več vozlišč z računskimi karticami Nvidia [V100](https://www.nvidia.com/en-us/data-center/v100/) in [H100](https://www.nvidia.com/en-us/data-center/h100/)

Primer zagona programa `nvidia-smi` (*nvidia system management interface*) na gruči. Program izpiše podatke o računskih karticah, ki so na voljo na danem računskem vozlišču.
```Bash
$ srun --partition=gpu --gpus=1 nvidia-smi --query
```

Napišimo še lasten [program](./koda/discover-device.cu) v programskem jeziku go, ki izpiše informacije o GPE. Podpora za programiranje grafičnih procesnih enot v programskem jeziku go je omejena. Uporabili bomo paket, ki v go doda možnost zaganjanja funkcij na grafičnih procesnih enotah s pomočjo okolja CUDA in ga je v okviru diplomske naloge razvil študent FRI. Paket najdete na [repozitoriju](https://github.com/InternatBlackhole/cudago). Na repozitoriju najdete tudi nekaj primerov uporabe paketa ter krajšo dokumentacijo. 

## Namestitev CudaGo na gruči Arnes

Najprej naložimo ustrezne module:
```Bash
$ module load CUDA
$ module load Go
```
Nato nastavimo okoljski spremenljivki `CGO_CFLAGS` in `CGO_LDFLAGS`, ki ju potrebujemo za namestitev prevajalnika CudaGo.
```Bash
export CGO_CFLAGS=$(pkg-config --cflags cudart-12.8) # or other version
export CGO_LDFLAGS=$(pkg-config --libs cudart-12.8) # or other version
```

Sedaj poženemo ukaz:
```Bash
$ go install github.com/InternatBlackhole/cudago/CudaGo@latest
```
S tem namestimo prevajalnik CudaGo v mapo `~/go/bin`. Da se izognemo pisanju polne poti, ko zaganjamo prevajalnik, dodamo v okoljsko spremenljivko `$PATH` ustrezno pot:
```Bash
$ export PATH="~/go/bin/:$PATH"
```
Preverimo, če `CudaGo` deluje:
```Bash
$ CudaGo -version
```

Okoljski spremenljivki `CGO_CFLAGS` in `CGO_LDFLAGS` in `$PATH` moramo nastaviti vsakič, ko želimo uporabljati prevajalnik CudaGo. V ta namen smo vam pripravili priročno [skripto](../../predavanja//21-cuda-primeri/koda/go/cudago-init.sh), ki ustrezno inicializira okolje.

## Izpis informacij o napravi CUDA v Go

Vzamemo kodo iz [primera](./koda/deviceInfo/main.go) in jo prenesemo v poljubno mapo na gruči.
Znotraj mape ustvarimo nov modul in namestimo potrebne pakete:
```Bash
$ go mod init cudaInfo
$ go mod tidy
```
Poženemo primer:
```Bash
$ srun --partition=gpu --gpus=1 go run .
```

## Zagon funkcije na GPE
Vzamemo datoteki [main.go](./koda/cudaHello/main.go) in [kernel.cu](./koda/cudaHello/kernel.cu) in ju prenesemo v poljubno mapo. Znotraj mape ustvarimo nov modul in namestimo potrebne pakete:
```Bash
$ go mod init cudaPrimer
$ go mod tidy
```

Sedaj prevedem ščepec znotraj datoteke `kernel.cu` v paket go `helloCuda`:
```Bash
$ CudaGo -package helloCuda kernel.cu
```

Poženemo primer:
```Bash
$ srun --partition=gpu --gpus=1 go run .
```

## Domača naloga 8

Vaša naloga je napisati program v Go, ki bo s pomočjo filtriranja z mediano odstranil šum iz slike. Rešitev (modul) v obliki datoteke zip oddajte preko [spletne učilnice]().

[Filitriranje z mediano](https://en.wikipedia.org/wiki/Median_filter) je postopek za odstranjevanje šuma v signalih. Pogosto se uporablja za predprocesiranje slik za pripravo za nadaljnje analize ali za izboljšanje vizualne kakovosti. Njegova prednost je, da ohranja robove in podrobnosti v slikah. Osnovna ideja je zamenjati vrednosti posamezne slikovne točke z mediano vrednosti točk v njeni okolici. Tipično se uporablja okno velikosti 3x3 ali 5x5 slikovnih točk.

Zašumljena slika             |  Slika po filtriranju z mediano
:-------------------------:|:-------------------------:
![](lenna-noisy.png)  |  ![](lenna-filtered.png)

## Postopek

* Filtriranje z mediano bomo zaradi enostavnosti izvajali nad sivinskimi slikami. Če slika ni sivinska jo naprej pretvorimo v sivinsko. Primer kako to naredimo v go, najdete [tukaj](./koda/imageCp/main.go).

* Izberemo velikost okna (3x3), ki določa območje sosednjih slikovnih točk okoli ciljne slikovne točke.

* Preberemo vse vrednosti slikovnih točk znotraj okna in jih uredimo po velikosti. 

* Mediano izračunamo tako, da vzamemo srednjo vrednost urejenega seznama. V primeru okna velikosti 3x3, imamo v seznamu 9 vrednosti. Srednja vrednost se nahaja na indeksu 4.

* Vrednost centralnega piksla nadomestimo z mediano.
  
![](postopek.png)          | 
:-------------------------:|
Postopek izračuna nove vrednosti ciljne slikovne točke   | 

* Postopek ponovimo za vse slikovne točke. 

V primeru robnih slikovnih točk, kjer nimamo na voljo vseh sosedov, obstaja več možnih pristopov. Najbolj običajen je, da za vrednosti slikovnih točk zunaj robov slike uporabimo kar najbližjo vrednost znotraj slike.

![](robovi.png)

## Branje in pisanje slik

Za branje in pisanje slik uporabite go-jev paket `image`. Primer uporabe paketa najdete [tukaj](./koda/imageCp/main.go).

## Filtriranje z mediano na GPE

Filtriranje izvedite na grafični procesni enoti s pomočjo ogrodja CUDA in paketa [CudaGo](https://github.com/InternatBlackhole/cudago). Pri tem je potrebno narediti naslednje korake:
 - Branje vhodne slike iz datoteke v pomnilnik gostitelja.
 - Pretvorba (če že ni) vhodne slike v sivinsko.
 - Rezervacija pomnilnika za podatkovne strukture na gostitelju (CPE) in napravi (GPE) (prostor za vhodno in izhodno sliko).
 - Prenos vhodne slike iz pomnilnika gostitelja na napravo.
 - Nastavitev organizacije niti: število blokov initi in število niti na blok. Uporabite 2D organizacijo niti, saj se najbolj prilega problemu, ki ga rešujemo. Pri nastavljanju organizacije niti moramu upoštevati tudi, kako bomo delo razdelili med niti. Uporabite pristop ena nit izračuna vrednost ene izhodne slikovne točke.
 - Zagon ščepca, ki izračuna izhodno sliko.
 - Prenos izhodne slike iz pomnilnika naprave v pomnilnik gostitelja.
 - Zapis izhodne slike v datoteko.

Za izhodišče lahko uporabite priloženo [kodo](./koda/imageCp/main.go). Pri iskanju mediane je potrebno urediti vrednosti slikovnih točk znotraj okna po velikosti. Poslužite se lahko poljubnega algortima za urejanje, pri tem upoštevajte, da je seznam vrednosti za urejanje majhen (9 števil). 

**Rok za oddajo: 13. 1. 2026**
