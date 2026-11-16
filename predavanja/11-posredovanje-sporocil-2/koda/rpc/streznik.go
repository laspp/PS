// Komunikacija po protokolu RPC
// strežnik

package main

import (
	"api/storage"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

func Server(url string, connHTTP bool) {
	// ustvarimo shrambo
	storage := storage.NewTodoStorage()
	// prijavimo metode za oddaljeno klicanje
	err := rpc.Register(storage)
	if err != nil {
		fmt.Println("Format of storage isn't correct. ", err)
	}
	// izpišemo ime strežnika
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	// odpremo vtičnico
	listener, err := net.Listen("tcp", url)
	if err != nil {
		panic(err)
	}
	fmt.Printf("RPC server listening at %v%v, HTTP=%v\n", hostName, url, connHTTP)
	// zaženemo strežnik
	if connHTTP {
		// pri HTTP pred streženjem prijavimo rokovalnik
		rpc.HandleHTTP()
		err = http.Serve(listener, nil)
		if err != nil {
			panic(err)
		}
	} else {
		rpc.Accept(listener)
	}
}
