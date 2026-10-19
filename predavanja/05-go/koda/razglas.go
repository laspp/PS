// sinhronizacija
// funkcija speaker da gorutinam listener signal kdaj lahko nadaljujejo
// s struct{} poudarimo, da kanal ni namenjen prenašanju sporočil

package main

import (
	"fmt"
	"time"
)

func speaker(message string) <-chan struct{} {
	broadcastStream := make(chan struct{})

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Announcement:", message)
		close(broadcastStream)
	}()

	return broadcastStream
}

func listener(id int, broadcastStream <-chan struct{}) {
	fmt.Println("Listener", id, "is waiting for an announcement.")

	<-broadcastStream

	fmt.Println("Listener", id, "completed.")
}

func main() {
	broadcastStream := speaker("Hello world for the last time!")

	for i := 1; i < 5; i++ {
		go listener(i, broadcastStream)
	}
	listener(0, broadcastStream)

	fmt.Println("Great!")

	time.Sleep(1 * time.Second)
}
