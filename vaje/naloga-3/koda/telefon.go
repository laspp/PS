package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/DistributedClocks/GoVector/govec"
)

type message struct {
	data   []byte
	length int
}

var start chan bool
var stopHeartbeat bool
var N int
var id int

var Logger *govec.GoLog
var opts govec.GoLogOptions

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
func receive(addr *net.UDPAddr) message {
	// Poslušamo
	conn, err := net.ListenUDP("udp", addr)
	checkError(err)
	defer conn.Close()
	fmt.Println("Telefon", id, "posluša na", addr)
	buffer := make([]byte, 1024)
	var msg []byte
	// Preberemo sporočilo
	_, err = conn.Read(buffer)
	checkError(err)
	Logger.UnpackReceive("Prejeto sporocilo ", buffer, &msg, opts)
	fmt.Println("Telefon", id, "prejel sporočilo:", string(msg))
	mLen := len(msg)
	// Vrnemo sporočilo
	rMsg := message{}
	rMsg.data = append(rMsg.data, msg[:mLen]...)
	rMsg.length = mLen
	return rMsg
}

func send(addr *net.UDPAddr, msg message) {
	// Odpremo povezavo
	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)
	defer conn.Close()
	// Pripravimo sporočilo
	Logger.LogLocalEvent("Priprava sporocila", opts)
	sMsg := fmt.Sprint(id) + "-"
	sMsg = string(msg.data[:msg.length]) + sMsg
	sMsgVC := Logger.PrepareSend("Poslano sporocilo ", []byte(sMsg), opts)
	_, err = conn.Write(sMsgVC)
	checkError(err)
	fmt.Println("Telefon", id, "poslal sporočilo", sMsg, "telefonu na naslovu", addr)

	// Ustavimo heartbeat servis
	stopHeartbeat = true
}

func heartBeat(addr *net.UDPAddr) {

	if id != 0 {
		// Ostali javljajo procesu 0, da so živi
		conn, err := net.DialUDP("udp", nil, addr)
		checkError(err)
		defer conn.Close()
		beat := [1]byte{byte(id)}
		for !stopHeartbeat {
			_, err = conn.Write(beat[:])
			time.Sleep(time.Second)
		}
	} else {
		// Posluša samo 0
		conn, err := net.ListenUDP("udp", addr)
		checkError(err)
		defer conn.Close()
		beat := make([]byte, 1)
		clients := make(map[byte]bool)
		for !stopHeartbeat {
			_, err = conn.Read(beat)
			checkError(err)
			fmt.Println("Telefon", id, "prejel utrip:", beat[:], len(clients))
			clients[beat[0]] = true
			// Če so se vsi javili zaključimo
			if len(clients) == N-1 {
				start <- true
				return
			}
		}
	}
}

func main() {
	// Preberi argumente
	portPtr := flag.Int("p", 9000, "# start port")
	idPtr := flag.Int("id", 0, "# process id")
	NPtr := flag.Int("n", 2, "total number of processes")
	flag.Parse()
	id = *idPtr
	// dnevnik z vektorsko uro
	Logger = govec.InitGoVector("Telefon-"+strconv.Itoa(id), "Log-Telefon-"+strconv.Itoa(id), govec.GetDefaultConfig())
	opts = govec.GetDefaultLogOptions()

	rootPort := *portPtr

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

	// Ustvari kanal, ki bo signaliziral, da so vsi procesi pripravljeni
	start = make(chan bool)

	// Zaženemo heartbeat servis, ki čaka, na javljanje vseh udeleženih procesov
	stopHeartbeat = false
	go heartBeat(rootAddr)

	// Izmenjava sporočil
	if id == 0 {
		<-start
		send(remoteAddr, message{})
		rMsg := receive(localAddr)
		fmt.Println(string(rMsg.data[:rMsg.length]) + "0")
	} else {
		rMsg := receive(localAddr)
		send(remoteAddr, rMsg)
	}

}
