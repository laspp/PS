// Komunikacija po protokolu HTTP (REST)
// odjemalec

package main

import (
	"api/storage"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func Client(url string) {

	// strukture, ki jih uporabljamo kot argumente pri klicu oddaljenih metod
	lecturesCreate := storage.Todo{Task: "predavanja", Completed: false}
	lecturesUpdate := storage.Todo{Task: "predavanja", Completed: true}
	practicals := storage.Todo{Task: "vaje", Completed: false}
	readAll := storage.Todo{Task: "", Completed: false}

	fmt.Printf("REST client connecting to %v\n", url)

	// ustvarimo zapis, varianta 1
	fmt.Print("1. Create: post    : ")
	if err := createWithPost(url, lecturesCreate); err == nil {
		fmt.Println("done")
	} else {
		panic(err)
	}

	// preberemo zapise, varianta 1
	fmt.Print("2. Read 1: get     : ")
	todoList := make(map[string]storage.Todo)
	if err := readWithGet(url, lecturesUpdate, &todoList); err == nil {
		fmt.Println(todoList, ": done")
	} else {
		panic(err)
	}

	// ustvarimo zapis, varianta 2
	fmt.Print("3. Create: req.post: ")
	if err := create(url, practicals); err == nil {
		fmt.Println("done")
	} else {
		panic(err)
	}

	// preberemo zapise, varianta 2
	fmt.Print("4. Read *: req.get : ")
	todoList = make(map[string]storage.Todo)
	if err := read(url, readAll, &todoList); err == nil {
		fmt.Println(todoList, ": done")
	} else {
		panic(err)
	}

	// posodobimo zapis
	fmt.Print("5. Update: req.put : ")
	if err := update(url, lecturesUpdate); err == nil {
		fmt.Println("done")
	} else {
		panic(err)
	}

	// izbrišemo zapis
	fmt.Print("6. Delete: req.del : ")
	if err := remove(url, practicals); err == nil {
		fmt.Println("done")
	} else {
		panic(err)
	}

	// še enkrat preberemo
	fmt.Print("7. Read *: req.get : ")
	todoList = make(map[string]storage.Todo)
	if err := read(url, readAll, &todoList); err == nil {
		fmt.Println(todoList, ": done")
	} else {
		panic(err)
	}
}

// ustvarimo zapis
func createWithPost(url string, todo storage.Todo) error {
	jsonBody, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close() // sprosti povezavo
	if resp.StatusCode != 200 {
		return errors.New("create not successful")
	}
	return nil
}

// preberemo zapise
func readWithGet(url string, todo storage.Todo, todoList *map[string]storage.Todo) error {
	resp, err := http.Get(url + todo.Task)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("read not successful")
	}
	err = json.NewDecoder(resp.Body).Decode(todoList)
	return err
}

// ustvarimo zapis
func create(url string, todo storage.Todo) error {
	jsonBody, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("create not successful")
	}
	return nil
}

// preberemo zapise
func read(url string, todo storage.Todo, todoList *map[string]storage.Todo) error {
	jsonBody, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("create not successful")
	}
	err = json.NewDecoder(resp.Body).Decode(todoList)
	return err
}

// posodobimo zapis
func update(url string, todo storage.Todo) error {
	jsonBody, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, url + todo.Task, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = errors.New("update not successful")
	}
	return err
}

// izbrišemo zapis
func remove(url string, todo storage.Todo) error {
	req, err := http.NewRequest(http.MethodDelete, url + todo.Task, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = errors.New("delete not successful")
	}
	return err
}
