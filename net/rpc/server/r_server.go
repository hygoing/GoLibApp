package main

import (
	"errors"
	"fmt"
	"net/rpc"
	"net"
	"net/http"
)

var c = make(chan string,1)

type Args struct {
	A, B int
}
type Quotient struct {
	Quo, Rem int
}
type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	fmt.Println("dasdasdasdasdas")
	c <- "111"
	fmt.Println("---%s---", c)
	*reply = args.A * args.B
	return nil
}
func (t *Arith) Divide(args *Args, quo *Quotient) error {
	fmt.Println(<-c)
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {
	ls, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println(err)
	}

	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()

	http.Serve(ls, nil)

}
