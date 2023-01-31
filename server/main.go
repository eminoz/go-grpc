package main

import (
	"fmt"
	"net"

	"github.com/eminoz/grpc-api/controller"
	api "github.com/eminoz/grpc-api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listen, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	newServer := grpc.NewServer()
	reflection.Register(newServer)
	api.RegisterEchoServer(newServer, &controller.Strm{})

	err = newServer.Serve(listen)
	if err != nil {
		panic(err)
	}
	fmt.Print("server stared on port 4040 ")
}
