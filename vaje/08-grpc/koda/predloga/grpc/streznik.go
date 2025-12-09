// Komunikacija po protokolu gRPC
// strežnik

package main

import (
	"api/grpc/protobufStorage"
	"api/storage"
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Server(url string) {
	// pripravimo strežnik gRPC
	grpcServer := grpc.NewServer()

	// pripravimo strukuro za streženje metod CRUD na shrambi TodoStorage
	crudServer := NewServerCRUD()

	// streženje metod CRUD na shrambi TodoStorage povežemo s strežnikom gRPC
	protobufStorage.RegisterCRUDServer(grpcServer, crudServer)

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
	fmt.Printf("gRPC server listening at %v%v\n", hostName, url)
	// začnemo s streženjem
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

// stuktura za strežnik CRUD za shrambo TodoStorage
type serverCRUD struct {
	protobufStorage.UnimplementedCRUDServer
	todoStore storage.TodoStorage
}

// pripravimo nov strežnik CRUD za shrambo TodoStorage
func NewServerCRUD() *serverCRUD {
	todoStorePtr := storage.NewTodoStorage()
	return &serverCRUD{protobufStorage.UnimplementedCRUDServer{}, *todoStorePtr}
}

// metode strežnika CRUD za shrambo TodoStorage
// odtis metode (argumenti, vrnjene vrednosti) so podane tako, kot zahtevajo paketi, ki jih je ustvaril prevajalnik protoc
// vse metode so samo ovojnice za klic metod v TodoStorage
// metode se izvajajo v okviru omejitev izvajalnega okolju ctx
func (s *serverCRUD) Create(ctx context.Context, in *protobufStorage.Todo) (*emptypb.Empty, error) {
	var ret struct{}
	err := s.todoStore.Create(&storage.Todo{Task: in.Task, Completed: in.Completed}, &ret)
	return &emptypb.Empty{}, err
}

func (s *serverCRUD) Read(ctx context.Context, in *protobufStorage.Todo) (*protobufStorage.TodoStorage, error) {
	dict := make(map[string](storage.Todo))
	err := s.todoStore.Read(&storage.Todo{Task: in.Task, Completed: in.Completed}, &dict)
	pbDict := protobufStorage.TodoStorage{}
	for k, v := range dict {
		pbDict.Todos = append(pbDict.Todos, &protobufStorage.Todo{Task: k, Completed: v.Completed})
	}
	return &pbDict, err
}

func (s *serverCRUD) Update(ctx context.Context, in *protobufStorage.Todo) (*emptypb.Empty, error) {
	var ret struct{}
	err := s.todoStore.Update(&storage.Todo{Task: in.Task, Completed: in.Completed}, &ret)
	return &emptypb.Empty{}, err
}

func (s *serverCRUD) Delete(ctx context.Context, in *protobufStorage.Todo) (*emptypb.Empty, error) {
	var ret struct{}
	err := s.todoStore.Delete(&storage.Todo{Task: in.Task, Completed: in.Completed}, &ret)
	return &emptypb.Empty{}, err
}
