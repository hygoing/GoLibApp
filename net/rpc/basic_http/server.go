package main

import (
	"fmt"
	"net/rpc"
	"net"
	"net/http"
	"GoLibApp/net/common"
)

func main() {
	ls, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println(err)
	}

	arith := new(common.Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()

	http.Serve(ls, nil)
}
