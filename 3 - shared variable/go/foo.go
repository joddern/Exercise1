// Use `go run foo.go` to run your program

package main

import (
    . "fmt"
    "runtime"
    "time"
)

var i = 0

// Want things (That are happening to the shared variable i) to happen sequentially, 
// not functions running after eachother sequentially, they will run simultanously.

// The two functions will take in channel where they can send 1

func incrementing(ch chan int, ending_ch chan int) {
    //TODO: increment i 1000000 times
    for k := 0 ; k < 1000000 ; k++ {
        ch <- 1
    }
    ending_ch <- 1
}

func decrementing(ch chan int, ending_ch chan int) {
    //TODO: decrement i 1000000 times
    for k := 0 ; k < 1000001 ; k++ {
        ch <- 1
    }
    ending_ch <- 1
}

func server(ch1, ch2, stop_ch chan int) {
    for {
        select {
            case <-ch1:
                i++
            case <-ch2:
                i--
            case <-stop_ch:
                return

        }
    }
}

func main() {
    // What does GOMAXPROCS do? What happens if you set it to 1?
    runtime.GOMAXPROCS(2)    

    // NEW STUFF START

    ch1 := make(chan int)
    ch2 := make(chan int)
    ending_channel := make(chan int)
    stop_channel := make(chan int)

    go incrementing(ch1, ending_channel)
    go decrementing(ch2, ending_channel)

    go func() {
        counter := 0
        for val := range ending_channel {
            counter += val
            if counter == 2 {
                stop_channel <- 1
            }
        }
    }()
    
    // Clearly crucial to have server not run as a goroutine, and call it after starting the other routines
    server(ch1, ch2, stop_channel)

    // NEW STUFF END
	
    // TODO: Spawn both functions as goroutines
    // go incrementing(inc_dec_ch)
    // go decrementing(inc_dec_ch)
	
    // We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
    // We will do it properly with channels soon. For now: Sleep.
    time.Sleep(500*time.Millisecond)
    Println("The magic number is:", i)
}
