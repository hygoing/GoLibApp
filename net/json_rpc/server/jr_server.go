package main

import (
	"net"
	"fmt"
	"net/rpc/jsonrpc"
	"log"
	"net/rpc"
	"errors"
)

type Args struct {
	A, B int
}
type Quotient struct {
	Quo, Rem int
}
type Arith int
func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}
func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {
	ls,err := net.Listen("tcp",":1234")
	if err != nil{
		fmt.Println(err)
	}
	defer ls.Close()

	arith := new(Arith)
	rpc.Register(arith)

	for {
		conn, err := ls.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
		}
		jsonrpc.ServeConn(conn)
		fmt.Println("sdsadas")
	}
}
