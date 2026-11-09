// Komunikacija po protokolu TCP
// 		niz pred pošiljanjem pretvorimo v []byte, ob prejemu pa []byte v niz
// odjemalec

package main

import (
	"fmt"
	"net"
	"time"
)

func Client(addr string, message string) {

	// povežemo se na strežnik
	connection, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	fmt.Println("TCP (string) client connected to", addr)

	// pošljemo sporočilo
	timeNow := time.Now().Format(time.DateTime)
	msgSend := fmt.Sprintf("%v @ %v", message, timeNow)
	fmt.Println("Sent message:", msgSend)
	_, err = connection.Write([]byte(msgSend))
	if err != nil {
		panic(err)
	}

	// sprejmemo odgovor
	// predpostavimo, da je odgovor krajši od 1024 bajtov in da prispe v enem kosu
	bMsgRecv := make([]byte, 1024)
	var msgLen int
	msgLen, err = connection.Read(bMsgRecv)
	if err != nil {
		panic(err)
	}
	fmt.Println("Received message:", string(bMsgRecv[0:msgLen]))
}
