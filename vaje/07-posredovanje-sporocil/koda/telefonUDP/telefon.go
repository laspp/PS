package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

type message struct {
	data   []byte
	length int
}

var start chan struct{}
var stopHeartbeat chan struct{}
var N int
var id int

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
func receive(conn *net.UDPConn) message {
	buffer := make([]byte, 1024)
	// Preberemo sporočilo
	mLen, err := conn.Read(buffer)
	checkError(err)
	fmt.Println("Telefon", id, "prejel sporočilo:", string(buffer[:mLen]))
	// Vrnemo sporočilo
	rMsg := message{}
	rMsg.data = append(rMsg.data, buffer[:mLen]...)
	rMsg.length = mLen
	return rMsg
}

func send(addr *net.UDPAddr, msg message) {
	// Odpremo povezavo
	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)
	defer conn.Close()
	// Pripravimo sporočilo
	sMsg := fmt.Sprint(id) + "-"
	sMsg = string(msg.data[:msg.length]) + sMsg
	_, err = conn.Write([]byte(sMsg))
	checkError(err)
	fmt.Println("Telefon", id, "poslal sporočilo", sMsg, "telefonu na naslovu", addr)
	// Ustavimo heartbeat servis
	close(stopHeartbeat)
}

func heartBeat(addr *net.UDPAddr) {

	if id != 0 {
		// Ostali javljajo procesu 0, da so pripravljeni
		conn, err := net.DialUDP("udp", nil, addr)
		checkError(err)
		defer conn.Close()
		beat := [1]byte{byte(id)}
		for {
			select {
			case <-stopHeartbeat:
				// Kanal je zaprt, prenehamo pošiljati utripe
				return
			default:
				conn.Write(beat[:])
				time.Sleep(time.Second)
			}
		}
	} else {
		// Posluša samo 0
		conn, err := net.ListenUDP("udp", addr)
		checkError(err)
		defer conn.Close()
		beat := make([]byte, 1)
		clients := make(map[byte]bool)
		// Čakamo, da se javijo vsi ostali procesi
		for len(clients) < N-1 {
			_, err = conn.Read(beat)
			checkError(err)
			clients[beat[0]] = true
			fmt.Println("Telefon", id, "prejel utrip od", beat[0], "- javljenih:", len(clients))
		}
		// Vsi so se javili, glavni proces lahko začne
		close(start)
	}
}

func main() {
	// Preberi argumente
	portPtr := flag.Int("p", 9000, "# start port")
	idPtr := flag.Int("id", 0, "# process id")
	NPtr := flag.Int("n", 2, "total number of processes")
	flag.Parse()

	rootPort := *portPtr
	id = *idPtr
	N = *NPtr
	basePort := rootPort + 1 + id
	nextPort := rootPort + 1 + ((id + 1) % N)

	// Ustvari potrebne mrežne naslove
	rootAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", rootPort))
	checkError(err)

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", basePort))
	checkError(err)

	remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", nextPort))
	checkError(err)

	// Ustvari kanala: start signalizira, da so vsi procesi pripravljeni,
	// z zaprtjem stopHeartbeat pa ustavimo pošiljanje utripov
	start = make(chan struct{})
	stopHeartbeat = make(chan struct{})

	// Najprej začnemo poslušati
	conn, err := net.ListenUDP("udp", localAddr)
	checkError(err)
	defer conn.Close()
	fmt.Println("Telefon", id, "posluša na", localAddr)

	// Zaženemo heartbeat servis, ki čaka na javljanje vseh udeleženih procesov
	go heartBeat(rootAddr)

	// Izmenjava sporočil
	if id == 0 {
		<-start
		send(remoteAddr, message{})
		rMsg := receive(conn)
		fmt.Println(string(rMsg.data[:rMsg.length]) + "0")
	} else {
		rMsg := receive(conn)
		send(remoteAddr, rMsg)
	}
}
