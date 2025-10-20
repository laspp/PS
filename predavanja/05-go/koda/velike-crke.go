// pretvarjanje malih znakov v sporočilu v velike
// sporočilo razdelimo na znake (rune) in jih pošljemo v kanal
// znake (rune) poberemo iz kanala, male znake spremenimo v velike in sestavimo sporočilo
// gorutina, ki ustvari kanal, ga tudi zapre
// uporabimo anonimno funkcijo, ki se izvaja v svoji gorutini, tudi potem, ko funkcija getLettersFromMessage konča

package main

import (
	"fmt"
	"unicode"
)

func getLettersFromMessage(message string) <-chan rune {

	defer fmt.Println("getLettersFromMessage: done")
	fmt.Println("getLettersFromMessage: start")
	letterStream := make(chan rune)

	go func() {
		defer fmt.Println("anonymous function: done") // prvi napisani defer se izvede zadnji
		defer close(letterStream)
		fmt.Println("anonymous function: start") // prvi napisani defer se izvede zadnji
		for _, letter := range message {
			letterStream <- letter
		}
	}()

	return letterStream
}

func getMessageFromLetters(letterStream <-chan rune) string {
	capsMessage := ""
	for letter := range letterStream {
		capsMessage += string(unicode.ToUpper(letter))
	}
	return capsMessage
}

func main() {
	var message string = "Hello world!"

	letterStream := getLettersFromMessage(message)

	capsMessage := getMessageFromLetters(letterStream)
	fmt.Println(message, " --> ", capsMessage)
}
