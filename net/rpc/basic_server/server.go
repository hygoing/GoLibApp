package main

import (
	"net/rpc"
	"GoLibApp/net/common"
	"net"
	"fmt"
)

/**
	rpc.Accept server.Accept 可监听每一个connect并serveConn
 */

func main() {
	ls, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Printf(err.Error())
	}
	server := rpc.NewServer()
	arith := new(common.Arith)
	server.Register(arith)

	/*server.Accept(ls)*/
	for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Printf(err.Error())
		}
		server.ServeConn(conn)
	}
}
