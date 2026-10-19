# Uvod [IPP:1]

## Zakaj potrebujemo vzporedno in porazdeljeno procesiranje

- znanost:
  - realistične simulacije, ceneje od eksperimentov v naravi
  - vremenska napoved, iskanje novih zdravil, simulacije tokov
- podatkovno rudarjenje:
  - zbiranje ogromnih količin podatkov
  - iskanje zakonitosti v podatkih, preverjanje modelov
  - statistična analiza in umetna inteligenca
- poslovni svet:
  - hitrejše obdelave več podatkov predstavljajo konkurenčno prednost
  - stranke želijo učinkovite in robustne rešitve
  - računalništvo v oblaku
- računalništvo na robu (*angl.* edge computing)
  - paradigma Industrija 4.0, senzorska omrežja

## Razvoj računalniških sistemov

- von Neumannova arhitektura: CPE, pomnilnik, V/I naprave
- dolgo časa (do 2004) je zmogljivost računalniških sistemov naraščala z miniaturizacijo (releji, elektronke, vedno manjši in več tranzistorjev, zviševanje frekvence ure) in izboljševanjem strojne opreme
- veljalo je Dennardovo pravilo: izboljšanje tehnologije izdelovanja tranzistorjev poveča njihovo število in hitrost, količina odvedene toplote na enoto površine pa se ne spremeni - zmanjšanje tehnologije izdelovanja za faktor dva poveča število tranzistorjev za štirikrat, vsak tranzistor dela dvakrat hitreje (krajše povezave) ob nespremenjeni količini odvedene toplote (napetost in tok se ustrezno zmanjšata, zato tudi poraba energije na tranzistor)
- karakteristike procesorjev skozi čas ([graf](https://www.karlrupp.net/2018/02/42-years-of-microprocessor-trend-data/))
- Moorov zakon: število tranzistorjev v integriraem vezju se podvoji približno vsaki dve leti ([graf](https://en.wikipedia.org/wiki/Transistor_count#/media/File:Moore's_Law_Transistor_Count_1970-2020.png))
- težava z odvajanjem toplote pri več kot 130 W porabe na procesor (zračno hlajenje)
- pojavljati se začnejo večjedrne arhitekture (po 2004), ki so energijsko bolj učinkovite
- s povečevanjem števila procesorskih jeder se ohranja veljavnost Moorovega zakona

## Večjedrni procesorji so energijsko bolj učinkoviti od enojedrnih

- poraba procesorja je odvisna od frekvence ure (f), napajalne napetosti (U) in kapacitivnosti vezja (C)
- višja frekvenca, več preklopov, večja poraba energije
- če pretiravamo, se čipi pregrevajo in prihaja do nezanesljivega delovanja

    $$ P = U_0I=U_0\frac{de}{dt}=U_0C\frac{dU}{dt}$$

    $$ U = U_0\sin(2\pi f t) $$

    $$ \frac{dU}{dt} = 2\pi fU_0\cos(2\pi f t)$$

    $$ P = 2\pi C f U_0^2 \cos(2\pi f t)$$

    $$ P = k C f U_0^2 $$

- ker imamo dve procesorski jedri, lahko s polovično frekvenco ure opravimo približno enako količino dela
- dve procesorski jedri lahko delata pri nižji napetosti, imata pa višjo kapacitivnost

  <img src="slike/poraba-eno-dve-jedri.png" alt="poraba energije: eno in dve jedri" width="75%">

## Vzporedenje postane ključno

- ključne omejitve, ki so pripeljale do večjedrnih procesorjev:
  - omejitev porabe moči
  - omejitev vzporedenja v procesorju
    - cevovod
    - špekulativno izvajanje
    - superskalarni procesorji
  - omejitve prepustnosti pomnilnika
    - procesorji so mnogo hitrejši od pomnilnikov
    - pomaga predpomnilnik, vendar se na večjedrnih procesorjih stvari zapletejo
    - problemi z raztegljivostjo (skalabilnostjo)

- podpora za razvoj vzporedne programske kode
  - prevajalniki, ki bi znali naše zaporedne programe samostojno pretvoriti v vzporedno obliko, so v razvojni fazi in so slabi
  - knjižnice za podporo vzporedenju
  - novi programski jeziki z vgrajeno podporo za vzporedenje

- rešitev
  - da dobro izkoristimo nove arhitekture, moramo sami pisati vzporedne programe
  - novi procesorji imajo mehanizme, ki znajo izkoriščati vzporedne programe (strojne niti)
  - za doseganje dobrih rezultatov moramo dobro poznati arhitekturo
  - kako pisati kodo, da bo učinkovita tudi na prihajajočih procesorjih?

## Aktualnost področja

- bolj učinkovito izvajanje algoritmov (hitrejše, manjša poraba energije)
- večja zanesljivost
- programska oprema kot množica mikrostoritev
- superračunalniški sistemi v Evropi in svetu ([www.top500.org](https://www.top500.org/))
- prosto dostopne superračunalniške gruče v Sloveniji ([https://doc.sling.si/navodila/clusters/](https://doc.sling.si/navodila/clusters/))

## Razlogi za delo na superračunalniški gruči

- vsi imamo enako strojno opremo (množica vozlišč z modernimi večjedrnimi procesorji in grafičnimi pospeševalniki)
- programsko opremo namestijo administratorji
- delate lahko doma ali na fakulteti
- celo študijsko leto imate dostop, lahko uporabljate tudi pri drugih predmetih, za podaljšanje pišete na <support@sling.si>
