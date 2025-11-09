// Komunikacija po protokolu TCP
//		s paketom gob pred pošiljanjem strukturo pretvorimo v []byte, ob sprejemu pa []byte v strukturo
// odjemalec

package main

import (
	"encoding/gob"
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
	fmt.Println("TCP (struct) client connected to", url)

	// pošljemo sporočilo
	msgSend := MessageAndTime{Message: message, Time: time.Now()}
	fmt.Println("Sent message:", msgSend)
	err = gob.NewEncoder(connection).Encode(&msgSend)
	if err != nil {
		panic(err)
	}

	// sprejmemo odgovor
	var msgRecv MessageAndTime
	err = gob.NewDecoder(connection).Decode(&msgRecv)
	if err != nil {
		panic(err)
	}
	fmt.Println("Received message:", msgRecv)
}
