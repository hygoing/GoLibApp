package main

import (
	"fmt"
	"os"
	"jd.com/jstack-common/log"
	"net/rpc"
	"GoLibApp/net/common"
)

func main() {
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := &common.Args{1,2}
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
}
