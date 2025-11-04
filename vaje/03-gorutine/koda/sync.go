/*
Program sync prikazuje uporabo sinhronizacijskih konstruktov v paketu sync
*/
package main

import (
	"fmt"
	"sync"
)

// Definiramo čakalno skupino
var wg sync.WaitGroup
var wc int

// Definiramo ključavnico
var lock sync.Mutex

// Delavec, ki povečuje števec
func workerInc(id int) {
	defer wg.Done()
	lock.Lock()
	wc++
	lock.Unlock()
}

// delavec, ki zmanjšuje števec
func workerDec(id int) {
	defer wg.Done()
	lock.Lock()
	wc--
	lock.Unlock()
}

func main() {
	workers := 100
	// Čakalno skupino inicializiramo z želenim številom delavcev
	wg.Add(2 * workers)
	// Zaženemo delavce
	for i := 0; i < workers; i++ {
		go workerInc(i)
		go workerDec(i)
	}
	// Počakamo, da delavci zaključijo
	wg.Wait()
	// Izpišemo končni rezultat
	// Kaj se zgodi, če iz delavcev odstranimo zaklepanje in odklepanje ključavnic?
	fmt.Println("Števec: ", wc)
}
