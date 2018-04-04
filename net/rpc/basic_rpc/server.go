package main

import (
	"net"
	"fmt"
	"GoLibApp/net/common"
	"net/rpc"
)

func main() {
	ls, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Printf(err.Error())
	}

	arith := new(common.Arith)
	rpc.Register(arith)

	for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Printf(err.Error())
		}
		go rpc.ServeConn(conn)
	}
}
