// stenski čas in monotoni čas
// v izpisu sta združena stenski čas in monotoni čas (m)
//		Time start: 2024-12-02 13:09:24.048793126 +0000 UTC m=+2.001038343
//		Time end  : 2024-12-02 13:09:25.048848009 +0000 UTC m=+3.001093326
//		Time elapsed (wall-clock): 1.000054883s
//		Time elapsed (monotonic) : 1.000054983s

package main

import (
	"fmt"
	"time"
)

func main() {

	time.Sleep(2 * time.Second)

	timeStart := time.Now()
	time.Sleep(1 * time.Second)
	timeEnd := time.Now()
	timeElapsed := timeEnd.Sub(timeStart)

	fmt.Printf("Time start: %v\n", timeStart)
	fmt.Printf("Time end  : %v\n", timeEnd)
	fmt.Printf("Time elapsed (wall-clock): %vs\n", float64(timeEnd.UnixNano()-timeStart.UnixNano())/1e9)
	fmt.Printf("Time elapsed (monotonic) : %v\n", timeElapsed)
}
