// Komunikacija po protokolu TCP
// 		niz pred pošiljanjem pretvorimo v []byte, ob prejemu pa []byte v niz
// odjemalec

package main

import (
	"fmt"
	"net"
	"time"
)

func Client(url string, message string) {

	// povežemo se na strežnik
	connection, err := net.Dial("tcp", url)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	fmt.Println("TCP (string) client connected to", url)

	// pošljemo sporočilo
	timeNow := time.Now().Format(time.DateTime)
	msgSend := fmt.Sprintf("%v @ %v", message, timeNow)
	fmt.Println("Sent message:", msgSend)
	_, err = connection.Write([]byte(msgSend))
	if err != nil {
		panic(err)
	}

	// sprejmemo odgovor
	bMsgRecv := make([]byte, 1024)
	var msgLen int
	msgLen, err = connection.Read(bMsgRecv)
	if err != nil {
		panic(err)
	}
	fmt.Println("Received message:", string(bMsgRecv[0:msgLen]))
}
