package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	netListen, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatal("Failed to listen %v", err.Error())
	}

	grpcServer := grpc.NewServer()

	fmt.Println("Server is now started at port: ", netListen.Addr())

	if err := grpcServer.Serve((netListen)); err != nil {
		log.Fatal("Failed to listen %v", err.Error())
	}
}