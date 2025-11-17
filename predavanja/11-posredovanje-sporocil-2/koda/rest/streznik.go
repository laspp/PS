// Komunikacija po protokolu HTTP (REST)
// strežnik

package main

import (
	"api/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Server(url string) {

	// ustvarimo shrambo in rokovalnik
	storage := storage.NewTodoStorage()
	storageHandler := NewTodosHandler(storage)

	// pripravimo http multiplekser
	// multiplekser glede na pot v url določi, kateri rokovalnik bo prevzel zahtevo
	mux := http.NewServeMux()
	// povežemo poti in rokovalnike
	mux.Handle("/", &HomeHandler{})
	mux.Handle("/todos", storageHandler)
	mux.Handle("/todos/", storageHandler)

	// preberemo ime vozlišča
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	// zaženemo strežnik
	fmt.Printf("REST server listening at %v%v\n", hostName, url)
	err = http.ListenAndServe(url, mux)
	if err != nil {
		panic(err)
	}
}

// rokovalnik za osnovno spletno stran
type HomeHandler struct{}

// strežnik za osnovno spletno stran
func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("RESTful CRUD server.\n\n"))
	w.Write([]byte("Usage:\n"))
	w.Write([]byte("\thttp://" + (*r).Host + "/todos/\n"))
	w.Write([]byte("\tREST method params: {task: <string>, completed: <bool>}\n"))
}

// rokovalnik za operacije CRUD (Create, Read, Update, Delete)
type TodosHandler struct {
	storage *storage.TodoStorage
}

// naredimo nov rokovalnik za naloge Todo
func NewTodosHandler(tds *storage.TodoStorage) *TodosHandler {
	return &TodosHandler{
		storage: tds,
	}
}

// strežnik za operacije CRUD
func (tdh *TodosHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		tdh.CreateTodo(w, r)
		return
	case r.Method == http.MethodGet:
		tdh.GetTodos(w, r)
		return
	case r.Method == http.MethodPut:
		tdh.UpdateTodo(w, r)
		return
	case r.Method == http.MethodDelete:
		tdh.DeleteTodo(w, r)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// metoda za obdelavo zahtevka za ustvarjanje zapisov
func (h *TodosHandler) CreateTodo(w http.ResponseWriter, r *http.Request) error {
	var todo storage.Todo
	var reply struct{}

	w.Header().Set("content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	if err := h.storage.Create(&todo, &reply); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// metoda za obdelavo zahtevka za branje zapisov
func (h *TodosHandler) GetTodos(w http.ResponseWriter, r *http.Request) error {
	var match string

	w.Header().Set("content-Type", "application/json")
	if strings.HasSuffix(r.URL.Path, "todos") || strings.HasSuffix(r.URL.Path, "todos/") {
		match = ""
	} else {
		split := strings.Split(r.URL.Path, "/")
		match = split[len(split)-1]
	}
	query := storage.Todo{Task: match, Completed: false}
	todoList := make(map[string]storage.Todo)
	err := h.storage.Read(&query, &todoList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	jsonBytes, err := json.Marshal(todoList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
	return nil
}

// metoda za obdelavo zahtevka za posodobitev zapisa
func (h *TodosHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) error {
	var todo storage.Todo
	var reply struct{}

	w.Header().Set("content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	if err := h.storage.Update(&todo, &reply); err != nil {
		if err == storage.ErrorNotFound {
			w.WriteHeader(http.StatusNotFound)
			return err
		}
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// metoda za obdelavo zahtevka za brisanje zapisa
func (h *TodosHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	var reply struct{}

	w.Header().Set("content-Type", "application/json")
	split := strings.Split(r.URL.Path, "/")
	match := split[len(split)-1]
	query := storage.Todo{Task: match, Completed: false}
	if err := h.storage.Delete(&query, &reply); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
