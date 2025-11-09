// Komunikacija po protokolu TCP
// 		niz pred pošiljanjem pretvorimo v []byte, ob prejemu pa []byte v niz
// strežnik

package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func Server(url string) {
	// izpišemo ime strežnika
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	fmt.Printf("TCP (string) server listening at %v%v\n", hostname, url)

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
	bMsgRecv := make([]byte, 1024)
	n, err := conn.Read(bMsgRecv)
	if err != nil {
		panic(err)
	}
	msgRecv := string(bMsgRecv[0:n])
	fmt.Println("Received message:", msgRecv)

	// obdelamo sporočilo
	time.Sleep(5 * time.Second)
	msgRecvSplit := strings.Split(msgRecv, "@")
	timeNow := time.Now().Format(time.DateTime)
	msgSend := fmt.Sprintf("Hello %v@ %v", msgRecvSplit[0], timeNow)

	// pošljemo odgovor
	fmt.Println("Sent message:", msgSend)
	_, err = conn.Write([]byte(msgSend))
	if err != nil {
		panic(err)
	}
}
