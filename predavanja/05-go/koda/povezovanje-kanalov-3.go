// povezovanje kanalov s stavkom select
// gorutini writer neprestano pišeta vsaka v svoj kanal
// gorutina reader bere in obdeluje sporočila iz več kanalov
// 		dodamo ključno besedo default - če ni sporočila, se vedno izvede ta, zato do prekoračitve časa pride izjemoma

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

func reader(stream1 <-chan string, stream2 <-chan string) {
	for {
		select {
		case msg1 := <-stream1:
			fmt.Println(msg1)
		case msg2 := <-stream2:
			fmt.Println(msg2)
		case <-time.After(1 * time.Second):
			fmt.Println("timeout")
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func main() {
	writer1Stream := writer(1)
	writer2Stream := writer(2)

	go reader(writer1Stream, writer2Stream)

	time.Sleep(20 * time.Second)
}
