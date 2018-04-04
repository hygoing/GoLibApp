package main

import (
	"net"
	"fmt"
	"net/rpc/jsonrpc"
	"log"
	"net/rpc"
	"GoLibApp/net/common"
)

func main() {
	ls, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println(err)
	}
	defer ls.Close()

	arith := new(common.Arith)
	rpc.Register(arith)

	for {
		conn, err := ls.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
		}
		go jsonrpc.ServeConn(conn)
	}
}
