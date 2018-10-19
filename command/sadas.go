package main

import (
	"time"
	"fmt"
)

type T struct {
	c chan string
}

func main() {
	t := &T{}
	go t.Init()
	time.Sleep(30 * time.Second)
	t.c = make(chan string, 100)

	time.Sleep(1000*time.Second)
}

func (t T) Init() {
	time.Sleep(5 * time.Second)
	t.c <- "aaa"
	fmt.Printf("sadasd")
}
