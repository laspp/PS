// povezovanje kanalov s stavkom select
// gorutini writer neprestano pišeta vsaka v svoj kanal
// gorutina reader bere in obdeluje sporočila iz več kanalov
//		dodamo kanal za zaustavitev gorutine reader
//		glavna gorutina sproži zaustavitev po 20 sekundah, gorutina reader zaključi in se pridruži
//		enako bi lahko naredili tudi za gorutini writer

package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

func writer(id int) <-chan string {
	stream := make(chan string)
	go func() {
		defer close(stream)
		for {
			stream <- "message from " + strconv.Itoa(id)
			time.Sleep(time.Duration(id*id+1) * time.Second)
		}
	}()
	return stream
}

func reader(stream1 <-chan string, stream2 <-chan string, streamDone <-chan struct{}) {
	defer wg.Done()
	for {
		select {
		case msg1 := <-stream1:
			fmt.Println(msg1)
		case msg2 := <-stream2:
			fmt.Println(msg2)
		case <-streamDone:
			fmt.Println("Done")
			return
		case <-time.After(1 * time.Second):
			fmt.Println("timeout")
		default:
		}
	}
}

func main() {
	writer1Stream := writer(1)
	writer2Stream := writer(2)
	streamDone := make(chan struct{})

	wg.Add(1)
	go reader(writer1Stream, writer2Stream, streamDone)

	time.Sleep(20 * time.Second)
	close(streamDone)

	wg.Wait()
}
