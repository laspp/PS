// primer uporabe paketa storage
// operacije CRUD kličemo lokalno

package main

import (
	"api/storage"
	"fmt"
)

func main() {
	var reply struct{}
	var replyMap map[string]storage.Todo

	// ustvarimo shrambo
	store := storage.NewTodoStorage()

	// pripravimo strukture, ki jih uporabljamo kot argumente
	lecturesCreate := storage.Todo{Task: "predavanja", Completed: false}
	lecturesUpdate := storage.Todo{Task: "predavanja", Completed: true}
	practicals := storage.Todo{Task: "vaje", Completed: false}
	readAll := storage.Todo{Task: "", Completed: false}

	// ustvarimo zapis
	fmt.Print("1. Create: ")
	if err := store.Create(&lecturesCreate, &reply); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// preberemo en zapis
	fmt.Print("2. Read 1: ")
	replyMap = make(map[string]storage.Todo)
	if err := store.Read(&lecturesUpdate, &replyMap); err != nil {
		panic(err)
	}
	fmt.Println(replyMap, ": done")

	// ustvarimo zapis
	fmt.Print("3. Create: ")
	if err := store.Create(&practicals, &reply); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// preberemo vse zapise
	fmt.Print("4. Read *: ")
	replyMap = make(map[string]storage.Todo)
	if err := store.Read(&readAll, &replyMap); err != nil {
		panic(err)
	}
	fmt.Println(replyMap, ": done")

	// posodobimo zapis
	fmt.Print("5. Update: ")
	if err := store.Update(&lecturesUpdate, &reply); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// izbrišemo zapis
	fmt.Print("6. Delete: ")
	if err := store.Delete(&practicals, &reply); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// preberemo vse zapise
	fmt.Print("7. Read *: ")
	replyMap = make(map[string]storage.Todo)
	if err := store.Read(&lecturesUpdate, &replyMap); err != nil {
		panic(err)
	}
	fmt.Println(replyMap, ": done")
}
