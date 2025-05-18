package main

import (
	"fmt"
	"time"
)

func hello(done chan bool) {
	fmt.Println("Hello go routine is going to sleep")
	time.Sleep(4 * time.Second)
	fmt.Println("Hello go routine awake and going to write to done")
	done <- true
}

func countdown() {
	for i := range [4]int{} {
		fmt.Println("countdown: ", i)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	done := make(chan bool)
	fmt.Println("Main going to call hello go goroutine")
	go hello(done)
	go countdown()
	<-done
	fmt.Println("Main recieved data")
}
