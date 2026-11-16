// Komunikacija po protokolu HTTP (REST)
// strežnik

package main

import (
	"api/storage"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

func Server(url string) {

	// ustvarimo shrambo in rokovalnik
	storage := storage.NewTodoStorage()
	storageHandler := NewTodosHandler(storage)

	// pripravimo http multiplekser
	// multiplekser glede na pot v url določi, kateri rokovalnik bo prevzel zahtevo
	// vzorec ima obliko "METODA /pot"; {task} je spremenljivka v poti,
	// {$} pomeni, da se pot konča natanko tu
	mux := http.NewServeMux()
	// povežemo poti in rokovalnike
	mux.Handle("GET /{$}", &HomeHandler{})
	mux.HandleFunc("POST /todos", storageHandler.CreateTodo)
	mux.HandleFunc("POST /todos/{$}", storageHandler.CreateTodo)
	mux.HandleFunc("GET /todos", storageHandler.GetTodos)
	mux.HandleFunc("GET /todos/{$}", storageHandler.GetTodos)
	mux.HandleFunc("GET /todos/{task}", storageHandler.GetTodos)
	mux.HandleFunc("PUT /todos/{task}", storageHandler.UpdateTodo)
	mux.HandleFunc("DELETE /todos/{task}", storageHandler.DeleteTodo)

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

// odgovor s kodo napake: 404, če naloge ni, sicer podana koda
func writeError(w http.ResponseWriter, err error, code int) {
	if errors.Is(err, storage.ErrorNotFound) {
		code = http.StatusNotFound
	}
	http.Error(w, err.Error(), code)
}

// metoda za obdelavo zahtevka za ustvarjanje zapisov
func (h *TodosHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo storage.Todo
	var reply struct{}

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if err := h.storage.Create(&todo, &reply); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// metoda za obdelavo zahtevka za branje zapisov
// pot /todos vrne vse zapise, pot /todos/{task} samo zapis z imenom task
func (h *TodosHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	query := storage.Todo{Task: r.PathValue("task")}
	todoList := make(map[string]storage.Todo)
	if err := h.storage.Read(&query, &todoList); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(todoList)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// metoda za obdelavo zahtevka za posodobitev zapisa
// ime naloge vzamemo iz poti, zato telo ne more spremeniti drugega zapisa
func (h *TodosHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var todo storage.Todo
	var reply struct{}

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	todo.Task = r.PathValue("task")
	if err := h.storage.Update(&todo, &reply); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// metoda za obdelavo zahtevka za brisanje zapisa
func (h *TodosHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	var reply struct{}

	query := storage.Todo{Task: r.PathValue("task")}
	if err := h.storage.Delete(&query, &reply); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
