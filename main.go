package main

import (
	"log"
)

func main() {
	var c chan int

	c = make(chan int, 1)

	c <- 10

	log.Println("print c:", <-c)

	c <- 11
	log.Println("print c:", <-c)
}
