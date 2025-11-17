// Komunikacija po protokolu RPC
// odjemalec

package main

import (
	"api/storage"
	"fmt"
	"net/rpc"
)

func Client(url string, connHTTP bool) {
	var client *rpc.Client
	var err error
	var reply struct{}
	var replyMap map[string]storage.Todo

	// strukture, ki jih uporabljamo kot argumente pri klicu oddaljenih metod
	lecturesCreate := storage.Todo{Task: "predavanja", Completed: false}
	lecturesUpdate := storage.Todo{Task: "predavanja", Completed: true}
	practicals := storage.Todo{Task: "vaje", Completed: false}
	readAll := storage.Todo{Task: "", Completed: false}

	// povežemo se na strežnik
	fmt.Printf("RPC client connecting to %v, HTTP=%v\n", url, connHTTP)
	if connHTTP {
		client, err = rpc.DialHTTP("tcp", url)
	} else {
		client, err = rpc.Dial("tcp", url)
	}
	if err != nil {
		panic(err)
	}

	// 1. ustvarimo zapis
	fmt.Print("1. Create: ")
	if err := client.Call("TodoStorage.Create", &lecturesCreate, &reply); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// 2. preberemo en zapis
	fmt.Print("2. Read 1: ")
	replyMap = make(map[string]storage.Todo)
	if err := client.Call("TodoStorage.Read", &lecturesUpdate, &replyMap); err != nil {
		panic(err)
	}
	fmt.Println(replyMap, ": done")

	// 3. ustvarimo zapis
	fmt.Print("3. Create: ")
	if err := client.Call("TodoStorage.Create", &practicals, &reply); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// 4. preberemo vse zapise
	fmt.Print("4. Read *: ")
	replyMap = make(map[string]storage.Todo)
	if err := client.Call("TodoStorage.Read", &readAll, &replyMap); err != nil {
		panic(err)
	}
	fmt.Println(replyMap, ": done")

	// 5. posodobimo zapis
	fmt.Print("5. Update: ")
	if err := client.Call("TodoStorage.Update", &lecturesUpdate, &reply); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// 6. izbrišemo zapis
	fmt.Print("6. Delete: ")
	if err := client.Call("TodoStorage.Delete", &practicals, &reply); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// 7. preberemo vse zapise (asinhrono)
	fmt.Print("7. Read *: ")
	replyMap = make(map[string]storage.Todo)
	rpcStream := client.Go("TodoStorage.Read", &lecturesUpdate, &replyMap, nil)
	rpcCall := <-rpcStream.Done
	if err := rpcCall.Error; err != nil {
		panic(err)
	}
	fmt.Println(replyMap, ": done")
}
