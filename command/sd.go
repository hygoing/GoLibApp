package main

import (
)

func main() {
	c := make(chan interface{})

	close(c)
	c <- "sdd"
}
