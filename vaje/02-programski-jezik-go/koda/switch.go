/*
Program switch prikazuje kako uporabljamo stavek switch v jeziku go
*/
package main

import (
	"fmt"
	"time"
)

func main() {

	// Tipičen primer stavka switch
	i := 2
	switch i {
	case 1:
		fmt.Println("ena")
	case 2:
		fmt.Println("dva")
	case 3:
		fmt.Println("tri")
	default:
		fmt.Println("nič od naštetega")
	}
	// Izpuščen izraz: ekvivalentno stavku if-else
	i = 3
	switch {
	case i == 1:
		fmt.Println("ena")
	case i == 2:
		fmt.Println("dva")
	case i == 3:
		fmt.Println("tri")
	default:
		fmt.Println("nič od naštetega")
	}
	// Primer inicializacije
	// Izraze lahko povežemo
	switch dayOfTheWeek := time.Now().Weekday(); dayOfTheWeek {
	case time.Saturday, time.Sunday:
		fmt.Println("Vikend je!")
	default:
		fmt.Println("Na FRI bo treba!")
	}
}
