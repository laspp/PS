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

var N int
var id int

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
func receive(addr *net.TCPAddr) message {
	// Poslušamo
	listener, err := net.ListenTCP("tcp", addr)
	checkError(err)
	fmt.Println("Telefon", id, "posluša na", addr)
	conn, err := listener.Accept()
	checkError(err)
	defer conn.Close()
	buffer := make([]byte, 1024)
	// Preberemo sporočilo
	mLen, err := conn.Read(buffer)
	checkError(err)
	fmt.Println("Telefon", id, "prejel:", string(buffer[:mLen]))
	// Vrnemo sporočilo
	rMsg := message{}
	rMsg.data = append(rMsg.data, buffer[:mLen]...)
	rMsg.length = mLen
	return rMsg
}

func send(addr *net.TCPAddr, msg message) {
	var err error
	retry := 5
	var conn *net.TCPConn
	// Odpremo povezavo, lahko je potrebnih več poskusov, če niso še vsi pripravljeni
	for retry != 0 {
		conn, err = net.DialTCP("tcp", nil, addr)
		retry--
		time.Sleep(time.Second / 2)
		if err == nil {
			break
		}
	}
	checkError(err)
	defer conn.Close()
	fmt.Println("Telefon", id, "se je povezal na", addr)
	// Pripravimo sporočilo
	sMsg := fmt.Sprint(id) + "-"
	sMsg = string(msg.data[:msg.length]) + sMsg
	_, err = conn.Write([]byte(sMsg))
	checkError(err)
	fmt.Println("Telefon", id, "poslal sporočilo", sMsg, "telefonu na naslovu", addr)

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
	basePort := rootPort + id
	nextPort := rootPort + ((id + 1) % N)
	// Ustvari potrebne mrežne naslove
	localAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", basePort))
	checkError(err)

	remoteAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", nextPort))
	checkError(err)

	// Izmenjava sporočil
	if id == 0 {
		send(remoteAddr, message{})
		rMsg := receive(localAddr)
		fmt.Println(string(rMsg.data[:rMsg.length]) + "0")
	} else {
		rMsg := receive(localAddr)
		send(remoteAddr, rMsg)
	}

}
