// Komunikacija po protokolu HTTP (REST)
// odjemalec

package main

import (
	"api/storage"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Client(url string) {
	fmt.Printf("REST client connecting to %v\n", url)

	// ustvarimo zapis, varianta 1
	fmt.Print("1. Create: post    : ")
	if err := createWithPost(url, []byte("{\"task\": \"predavanja\", \"completed\": false}")); err == nil {
		fmt.Println("done")
	} else {
		panic(err)
	}

	// preberemo zapise, varianta 1
	fmt.Print("2. Read 1: get     : ")
	if body, err := readWithGet(url, "predavanja"); err == nil {
		dict := make(map[string]storage.Todo)
		err = json.Unmarshal(body, &dict)
		if err != nil {
			panic(err)
		}
		fmt.Println(dict, ": done")
	} else {
		panic(err)
	}

	// ustvarimo zapis, varianta 2
	fmt.Print("3. Create: req.post: ")
	todo := storage.Todo{Task: "vaje", Completed: false}
	jsonBody, err := json.Marshal(&todo)
	if err != nil {
		panic(err)
	}
	if err := create(url, jsonBody); err == nil {
		fmt.Println("done")
	} else {
		panic(err)
	}

	// preberemo zapise, varianta 2
	fmt.Print("4. Read *: req.get : ")
	if body, err := read(url, ""); err == nil {
		dict := make(map[string]storage.Todo)
		err := json.Unmarshal(body, &dict)
		if err != nil {
			panic(err)
		}
		fmt.Println(dict, ": done")
	} else {
		panic(err)
	}

	// posodobimo zapis
	fmt.Print("5. Update: req.put : ")
	todo = storage.Todo{Task: "predavanja", Completed: true}
	jsonBody, err = json.Marshal(&todo)
	if err != nil {
		panic(err)
	}
	if err := update(url, "predavanja", jsonBody); err == nil {
		fmt.Println("done")
	} else {
		panic(err)
	}

	// izbrišemo zapis
	fmt.Print("6. Delete: req.del : ")
	if err := delete(url, "vaje"); err == nil {
		fmt.Println("done")
	} else {
		panic(err)
	}

	// še enkrat preberemo
	fmt.Print("7. Read *: req.get : ")
	if body, err := read(url, ""); err == nil {
		dict := make(map[string]storage.Todo)
		err = json.Unmarshal(body, &dict)
		if err != nil {
			panic(err)
		}
		fmt.Println(dict, ": done")
	} else {
		panic(err)
	}
}

// ustvarimo zapis
func createWithPost(url string, jsonBody []byte) error {
	bodyReader := bytes.NewReader(jsonBody)
	resp, err := http.Post(url, "application/json", bodyReader)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		err = errors.New("create not successful")
	}
	return err
}

// preberemo zapise
func readWithGet(url string, task string) ([]byte, error) {
	resp, err := http.Get(url + task)
	if err != nil {
		return []byte{}, err
	}
	resBody, err := io.ReadAll(resp.Body)
	return resBody, err
}

// ustvarimo zapis
func create(url string, body []byte) error {
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		err = errors.New("create not successful")
	}
	return err
}

// preberemo zapise
func read(url string, task string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url+task, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	body, err := io.ReadAll(resp.Body)
	return body, err
}

// posodobimo zapis
func update(url string, task string, body []byte) error {
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(http.MethodPut, url+task, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		err = errors.New("update not successful")
	}
	return err
}

// izbrišemo zapis
func delete(url string, task string) error {
	req, err := http.NewRequest(http.MethodDelete, url+task, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		err = errors.New("delete not successful")
	}
	return err
}
