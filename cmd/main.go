package main

import (
	"fmt"
	"log"
	"net"

	"go-grpc-unary/cmd/services"
	productPb "go-grpc-unary/pb/product"

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


	//gRPC SERVER
	grpcServer := grpc.NewServer()

	productService := services.ProductService{}
	productPb.RegisterProductServiceServer(grpcServer, &productService)

	fmt.Println("Server is now started at port: ", netListen.Addr())

	if err := grpcServer.Serve((netListen)); err != nil {
		log.Fatal("Failed to listen %v", err.Error())
	}

	
}