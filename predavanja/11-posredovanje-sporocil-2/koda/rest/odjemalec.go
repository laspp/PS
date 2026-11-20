// Komunikacija po protokolu HTTP (REST)
// odjemalec

package main

import (
	"api/storage"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Client(url string) {
	var reply struct{}
	var replyMap map[string]storage.Todo

	// pripravimo odjemalca REST
	fmt.Printf("REST client connecting to %v\n", url)
	restClient := NewTodoClient(url)

	// strukture, ki jih uporabljamo kot argumente pri klicu oddaljenih metod
	lecturesCreate := storage.Todo{Task: "predavanja", Completed: false}
	lecturesUpdate := storage.Todo{Task: "predavanja", Completed: true}
	practicals := storage.Todo{Task: "vaje", Completed: false}
	readAll := storage.Todo{Task: "", Completed: false}

	// 1. ustvarimo zapis
	fmt.Print("1. Create: ")
	if err := restClient.Create(&lecturesCreate, &reply); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// 2. preberemo en zapis
	fmt.Print("2. Read 1: ")
	replyMap = make(map[string]storage.Todo)
	if err := restClient.Read(&lecturesCreate, &replyMap); err != nil {
		panic(err)
	}
	fmt.Println(replyMap, ": done")

	// 3. ustvarimo zapis
	fmt.Print("3. Create: ")
	if err := restClient.CreateV2(&practicals, &reply); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// 4. preberemo vse zapise
	fmt.Print("4. Read *: ")
	replyMap = make(map[string]storage.Todo)
	if err := restClient.ReadV2(&readAll, &replyMap); err != nil {
		panic(err)
	}
	fmt.Println(replyMap, ": done")

	// 5. posodobimo zapis
	fmt.Print("5. Update: ")
	if err := restClient.Update(&lecturesUpdate, &reply); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// 6. izbrišemo zapis
	fmt.Print("6. Delete: ")
	if err := restClient.Delete(&practicals, &reply); err != nil {
		panic(err)
	}
	fmt.Println("done")

	// 7. preberemo vse zapise
	fmt.Print("7. Read *: ")
	replyMap = make(map[string]storage.Todo)
	if err := restClient.Read(&readAll, &replyMap); err != nil {
		panic(err)
	}
	fmt.Println(replyMap, ": done")
}

// odjemalec REST za shrambo nalog
//
//	metode imajo enak odtis kot metode storage.TodoStorage, zato jih kličemo
//	enako kot lokalno shrambo (storage/uporaba/uporaba.go);
//	vsak klic metode prevedemo v zahtevo HTTP, odgovor pa nazaj v vrnjene vrednosti
type TodoClient struct {
	baseURL string       	// na primer http://localhost:9876/todos/
	client  *http.Client 	// si zapomni odprte povezave in jih ponovno uporabi
}

// naredimo novega odjemalca za strežnik na naslovu baseURL
func NewTodoClient(baseURL string) *TodoClient {
	return &TodoClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

// ustvarimo zapis: POST /todos/ z nalogo v telesu
func (c *TodoClient) Create(todo *storage.Todo, ret *struct{}) error {
	body, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	resp, err := c.client.Post(c.baseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close() // sprostimo povezavo
	return checkStatus(resp)
}

// preberemo zapise: GET /todos/{task}, pri praznem imenu naloge GET /todos/ vrne vse zapise
func (c *TodoClient) Read(todo *storage.Todo, dict *map[string]storage.Todo) error {
	resp, err := c.client.Get(c.taskURL(todo))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return err
	}
	return json.NewDecoder(resp.Body).Decode(dict)
}

// ustvarimo zapis: POST /todos/ z nalogo v telesu
func (c *TodoClient) CreateV2(todo *storage.Todo, ret *struct{}) error {
	body, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json") // oblika telesa, ki ga pošiljamo
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

// preberemo zapise: GET /todos/{task}, pri praznem imenu naloge GET /todos/ vrne vse zapise
func (c *TodoClient) ReadV2(todo *storage.Todo, dict *map[string]storage.Todo) error {
	req, err := http.NewRequest(http.MethodGet, c.taskURL(todo), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json") // oblika odgovora, ki ga želimo
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := checkStatus(resp); err != nil {
		return err
	}
	return json.NewDecoder(resp.Body).Decode(dict)
}

// posodobimo zapis: PUT /todos/{task} z nalogo v telesu
func (c *TodoClient) Update(todo *storage.Todo, ret *struct{}) error {
	body, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, c.taskURL(todo), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

// izbrišemo zapis: DELETE /todos/{task}
func (c *TodoClient) Delete(todo *storage.Todo, ret *struct{}) error {
	req, err := http.NewRequest(http.MethodDelete, c.taskURL(todo), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkStatus(resp)
}

// naslov vira za nalogo; ime zakodiramo, da je veljaven del poti (presledek -> %20)
func (c *TodoClient) taskURL(todo *storage.Todo) string {
	return c.baseURL + url.PathEscape(todo.Task)
}

// statusno kodo, ki ni 200 OK, pretvorimo v napako s sporočilom strežnika
func checkStatus(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	msg, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("%v: %v", resp.Status, strings.TrimSpace(string(msg)))
}
