// Komunikacija po protokolu gRPC
// odjemalec

package main

import (
	"api/grpc/protobufStorage"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Client(url string) {
	// vzpostavimo povezavo s strežnikom
	fmt.Printf("gRPC client connecting to %v\n", url)
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// vzpostavimo izvajalno okolje
	contextCRUD, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// vzpostavimo vmesnik gRPC
	grpcClient := protobufStorage.NewCRUDClient(conn)

	// pripravimo strukture, ki jih uporabljamo kot argumente pri klicu oddaljenih metod
	lecturesCreate := protobufStorage.Todo{Task: "predavanja", Completed: false}
	lecturesUpdate := protobufStorage.Todo{Task: "predavanja", Completed: true}
	practicals := protobufStorage.Todo{Task: "vaje", Completed: false}
	readAll := protobufStorage.Todo{Task: "", Completed: false}

	// ustvarimo zapis
	fmt.Print("1. Create: ")
	if _, err := grpcClient.Create(contextCRUD, &lecturesCreate); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// preberemo en zapis
	fmt.Print("2. Read 1: ")
	if response, err := grpcClient.Read(contextCRUD, &lecturesCreate); err == nil {
		fmt.Println(response.Todos, ": done")
	} else {
		panic(err)
	}

	// ustvarimo zapis
	fmt.Print("3. Create: ")
	if _, err := grpcClient.Create(contextCRUD, &practicals); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// preberemo vse zapise
	fmt.Print("4. Read *: ")
	if response, err := grpcClient.Read(contextCRUD, &readAll); err == nil {
		fmt.Println(response.Todos, ": done")
	} else {
		panic(err)
	}

	// posodobimo zapis
	fmt.Print("5. Update: ")
	if _, err := grpcClient.Update(contextCRUD, &lecturesUpdate); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// izbrišemo zapis
	fmt.Print("6. Delete: ")
	if _, err := grpcClient.Delete(contextCRUD, &practicals); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// preberemo vse zapise
	fmt.Print("7. Read *: ")
	if response, err := grpcClient.Read(contextCRUD, &readAll); err == nil {
		fmt.Println(response.Todos, ": done")
	} else {
		panic(err)
	}
}
