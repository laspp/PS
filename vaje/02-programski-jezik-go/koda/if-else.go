/*
Program if-else prikazuje kako uporabljamo pogoj if v jeziku go
*/
package main

import "fmt"

func main() {
	// Tipičen primer stavka if-else
	// Zaviti oklepaji so obvezni
	if 3%2 == 0 {
		fmt.Println("3 je sodo število")
	} else {
		fmt.Println("3 je liho število")
	}

	// Stavek if-else if-else z inicializacijo
	if num := 42; num < 0 {
		fmt.Println(num, "je negativno število")
	} else if num > 0 {
		fmt.Println(num, "je pozitivno število")
	} else {
		fmt.Println(num, "je enako nič")
	}
}
