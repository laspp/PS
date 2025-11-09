// Komunikacija po protokolu TCP
//		s paketom gob pred pošiljanjem strukturo pretvorimo v []byte, ob sprejemu pa []byte v strukturo
// strežnik

package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"time"
)

func Server(url string) {
	// izpišemo ime strežnika
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	fmt.Printf("TCP (struct) server listening at %v%v\n", hostname, url)

	// odpremo vtičnico
	listener, err := net.Listen("tcp", url)
	if err != nil {
		panic(err)
	}
	for {
		// čakamo na odjemalca
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// obdelamo zahtevo
		go handleRequest(connection)
	}
}

func handleRequest(conn net.Conn) {
	// na koncu zapremo povezavo
	defer conn.Close()

	// sprejmemo sporočilo
	var msgRecv MessageAndTime
	err := gob.NewDecoder(conn).Decode(&msgRecv)
	if err != nil {
		panic(err)
	}
	fmt.Println("Received message:", msgRecv)

	// obdelamo sporočilo
	time.Sleep(5 * time.Second)
	msgSend := MessageAndTime{Message: "Hello " + msgRecv.Message, Time: time.Now()}

	// pošljemo odgovor
	fmt.Println("Sent message:", msgSend)
	err = gob.NewEncoder(conn).Encode(&msgSend)
	if err != nil {
		panic(err)
	}
}
