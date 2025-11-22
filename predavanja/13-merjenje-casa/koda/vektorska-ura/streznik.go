// Komunikacija po protokolu TCP z beleženjem z vektorsko uro
// 		niz pred pošiljanjem pretvorimo v []byte, ob prejemu pa []byte v niz
// strežnik

package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/DistributedClocks/GoVector/govec"
)

func Server(url string) {
	// dnevnik z vektorsko uro
	Logger := govec.InitGoVector("Server", "Log-Server", govec.GetDefaultConfig())
	opts := govec.GetDefaultLogOptions()

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
		go handleRequest(connection, Logger, opts)
	}
}

func handleRequest(conn net.Conn, logger *govec.GoLog, opts govec.GoLogOptions) {
	// na koncu zapremo povezavo
	defer conn.Close()

	// sprejmemo sporočilo
	bMsgRecvVC := make([]byte, 1024)
	var bMsgRecv []byte
	_, err := conn.Read(bMsgRecvVC)
	if err != nil {
		panic(err)
	}
	// sprejmemo sporočilo in iz njega odstranimo vektorsko uro
	logger.UnpackReceive("Received Message from client", bMsgRecvVC, &bMsgRecv, opts)
	msgRecv := string(bMsgRecv)
	fmt.Println("Received message:", msgRecv)

	// obdelamo sporočilo
	logger.LogLocalEvent("Server preparing reply", opts)
	time.Sleep(5 * time.Second)
	msgRecvSplit := strings.Split(msgRecv, "@")
	timeNow := time.Now().Format(time.DateTime)
	msgSend := fmt.Sprintf("Hello %v@ %v", msgRecvSplit[0], timeNow)

	// pošljemo odgovor
	fmt.Println("Sent message:", msgSend)
	// sporočilu dodamo vektorsko uro in ga pošljemo odjemalcu
	msgSendVC := logger.PrepareSend("Sending Message", []byte(msgSend), opts)
	_, err = conn.Write(msgSendVC)
	if err != nil {
		panic(err)
	}
}
