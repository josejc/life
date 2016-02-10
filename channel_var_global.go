package main

import "fmt"
import "time"

var c = make(chan int)
var a int

func f() {
	for {
		<-c
		fmt.Println(a)
	}
}

func main() {
	a = 5
	go f()
	a = 6
	c <- 0
	a = 7
	c <- 0
	time.Sleep(100 * time.Millisecond)
}
