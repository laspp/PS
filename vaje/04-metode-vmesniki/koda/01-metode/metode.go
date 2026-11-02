/*
Program prikazuje uporabo metod v programskem jeziku Go
*/
package main

import (
	"fmt"
)

// Definiramo strukturo Student
type Student struct {
	name    string
	surname string
	id      string
	year    int
}

// Funkcija za izpis imena in priimka študenta
func displayNameSurname(s Student) {
	fmt.Printf("Ime: %s, Priimek: %s\n", s.name, s.surname)
}

// Metoda za izpis imena in priimka študenta
func (s Student) displayNameSurname() {
	fmt.Printf("Ime: %s, Priimek: %s\n", s.name, s.surname)
}

// Metoda, ki vrača letnik študija
func (s Student) getYear() int {
	return s.year
}

// Metoda, ki nastavi letnik študija
// Uporabimo kazalec za prejemnika
func (s *Student) setYear(year int) {
	s.year = year
}

// Metoda, ki poskuša nastaviti letnik študija
// Prejemnik je vrednost, zato metoda spremeni le svojo kopijo študenta
func (s Student) setYearCopy(year int) {
	s.year = year
}

func main() {

	// Ustvarimo novega študenta
	student1 := Student{name: "Janez", surname: "Novak", id: "63230000", year: 1}

	// Pokličemo funkcijo
	displayNameSurname(student1)

	// Pokličemo metode
	student1.displayNameSurname()
	fmt.Println(student1.getYear())
	student1.setYear(2)
	fmt.Println(student1.getYear())
	// Letnik ostane nespremenjen, saj metoda spremeni le kopijo
	student1.setYearCopy(3)
	fmt.Println(student1.getYear())
}
