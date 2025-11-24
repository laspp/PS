// Protokol NTP
// 		protokol NTP meri čas v sekundah od 1.1.1900, čas Unix pa se meri v sekundah od 1.1.1970
//		v telegramu NTP so časi podani v
//			- sekundah po 1.1.1900 (...TimeSec) in
// 			- deležu sekunde (...TimeFrac / 2^32)

package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"time"
)

const ntpEpochOffset int64 = 2208988800 // sekunde med 1.1.1900 in 1.1.1970, upoštevanih je 17 prestopnih let

func timeNTPtoUnix(sec uint32, frac uint32) (timeUnix time.Time) {
	secs := int64(sec) - ntpEpochOffset
	nanos := (int64(frac) * 1e9) >> 32
	timeUnix = time.Unix(secs, nanos)
	return
}

func timeUnixToNTP(timeUnix time.Time) (sec uint32, frac uint32) {
	sec = uint32(timeUnix.Unix() + ntpEpochOffset)
	frac = uint32(((timeUnix.UnixNano() % 1e9) << 32) / 1e9)
	return
}

// struktura telegrama NTP v3, zahteva in odgovor uporabljata isti telegram
type telegram struct {
	Settings       uint8  // prestopna sekunda (2 bita), verzija (3 biti), režim (3 biti)
	Stratum        uint8  // oznaka strežnika, nižja številka pomeni natančnejši strežnik
	Poll           int8   // največji dovoljeni interval med zaporednimi sporočili
	Precision      int8   // natančnost strežnikove ure (2^exponent)
	RootDelay      uint32 // čas delta pri komunikaciji z nadrejenim strežnikom
	RootDispersion uint32 // največja napaka glede na primarni vir
	ReferenceID    uint32 // oznaka referenčne ure
	RefTimeSec     uint32 // čas zadnjega popravka ure na strežniku, sekunde
	RefTimeFrac    uint32 // čas zadnjega popravka ure na strežniku, del sekunde
	OrigTimeSec    uint32 // čas odjemalca ob pošiljanju zahtevka, sekunde
	OrigTimeFrac   uint32 // čas odjemalca ob pošiljanju zahtevka, del sekunde
	RxTimeSec      uint32 // čas sprejetja zahteve, sekunde
	RxTimeFrac     uint32 // čas sprejetja zahteve, del sekunde
	TxTimeSec      uint32 // čas odpošiljanja odgovora, sekunde
	TxTimeFrac     uint32 // čas odpošiljanja odgovora, del sekunde
}

func main() {
	// preberemo argumente
	sPtr := flag.String("s", "ntp1.arnes.si", "NTP server address")
	flag.Parse()

	// vzpostavimo povezavo UDP s strežnikom
	conn, err := net.Dial("udp", *sPtr+":123")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// za branje ali pošiljanje čakamo največ 3 sekunde
	if err := conn.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
		panic(err)
	}

	// pošljemo zahtevo na strežnik NTP
	// Settings: 0x1B = 00.011.011 = informacija o prestopni sekundi, verzija 3, tip zahteve
	// tip zahteve: (1 - brat-bratu-aktivni, 2 - brat-bratu-pasivni, 3 - odjemalec, 4 - strežnik, 5 - razširjanje, ...)
	time1 := time.Now()
	secs1, nanos1 := timeUnixToNTP(time1)
	req := &telegram{Settings: 0x1B, TxTimeSec: secs1, TxTimeFrac: nanos1}
	if err := binary.Write(conn, binary.BigEndian, req); err != nil {
		panic(err)
	}

	// počakamo odgovor
	rsp := &telegram{}
	if err := binary.Read(conn, binary.BigEndian, rsp); err != nil {
		panic(err)
	}
	// zabeležimo čas, ko smo odgovor prejeli
	time4 := time.Now()
	// preberemo vsebino odgovora
	timeRef := timeNTPtoUnix(rsp.RefTimeSec, rsp.RefTimeFrac)
	time1telegram := timeNTPtoUnix(rsp.OrigTimeSec, rsp.OrigTimeFrac)
	time2 := timeNTPtoUnix(rsp.RxTimeSec, rsp.RxTimeFrac)
	time3 := timeNTPtoUnix(rsp.TxTimeSec, rsp.TxTimeFrac)

	// izračunamo čas prenosa tja in nazaj
	delta := (time4.UnixNano() - time1telegram.UnixNano()) - (time3.UnixNano() - time2.UnixNano())
	// ocenimo razliko med časom strežnika in odjemalca
	// theta > 0: odjemalčeva ura zaostaja
	theta := (time3.UnixNano() + delta/2) - time4.UnixNano()

	fmt.Printf("Server: %v\n", *sPtr)
	fmt.Printf("Telegram (req): %#v\n", req)
	fmt.Printf("Telegram (res): %#v\n", rsp)
	fmt.Printf("Tref : %v\n", timeRef)
	fmt.Printf("T1   : %v\n", time1)
	fmt.Printf("T1tel: %v\n", time1telegram)
	fmt.Printf("T2   : %v\n", time2)
	fmt.Printf("T3   : %v\n", time3)
	fmt.Printf("T4   : %v\n", time4)
	fmt.Printf("delta: %v ns = %v s\n", delta, float64(delta)/1e9)
	fmt.Printf("theta: %v ns = %v s\n", theta, float64(theta)/1e9)
}
