// Komunikacija po protokolu TCP z beleženjem z vektorsko uro
// 		niz pred pošiljanjem pretvorimo v []byte, ob prejemu pa []byte v niz
// odjemalec

package main

import (
	"fmt"
	"net"
	"time"

	"github.com/DistributedClocks/GoVector/govec"
)

func Client(url string, message string) {
	// dnevnik z vektorsko uro
	Logger := govec.InitGoVector("Client", "Log-Client", govec.GetDefaultConfig())
	opts := govec.GetDefaultLogOptions()

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
	msgSendVC := Logger.PrepareSend("Sending Message", []byte(msgSend), opts)

	_, err = connection.Write(msgSendVC)
	if err != nil {
		panic(err)
	}

	// sprejmemo odgovor
	bMsgRecvVC := make([]byte, 1024)
	var bMsgRecv []byte
	_, err = connection.Read(bMsgRecvVC)
	if err != nil {
		panic(err)
	}
	Logger.UnpackReceive("Received Message from server", bMsgRecvVC, &bMsgRecv, opts)
	fmt.Println("Received message:", string(bMsgRecv))

	Logger.LogLocalEvent("Client finished", opts)
}
