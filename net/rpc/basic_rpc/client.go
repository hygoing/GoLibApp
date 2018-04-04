package main

import (
	"fmt"
	"net/rpc"
	"os"
	"GoLibApp/net/common"
)

func main() {
	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		fmt.Printf(err.Error())
	}

	args := &common.Args{1, 2}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		fmt.Println("Arith.Muliply call error:", err)
		os.Exit(1)
	}
	fmt.Println("the arith.mutiply is :", reply)

	var quto common.Quotient
	err = client.Call("Arith.Divide", args, &quto)
	if err != nil {
		fmt.Println("Arith.Divide call error:", err)
		os.Exit(1)
	}
	fmt.Println("the arith.devide is :", quto.Quo, quto.Rem)

	divCall := client.Go("Arith.Multiply", args, &reply, nil)
	replyCall := <-divCall.Done
	fmt.Printf("%s,----------,%s", divCall, replyCall)
}
