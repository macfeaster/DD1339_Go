package main

import (
	"fmt"
	"time"
)

// This program should go to 11, but sometimes it only prints 1 to 10.
func main() {
	ch := make(chan int)
	done := make(chan bool)
	go Print(ch, done)
	for i := 1; i <= 11; i++ {
		ch <- i
	}
	close(ch)
	<-done
}

// Print prints all numbers sent on the channel.
// The function returns when the channel is closed.
func Print(ch <-chan int, done chan bool) {
	for n := range ch { // reads from channel until it's closed
		time.Sleep(time.Millisecond * 250)
		fmt.Println(n)
	}
	done <- true
}
