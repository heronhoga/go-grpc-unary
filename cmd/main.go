package main

import (
	"fmt"
	"go-grpc-unary/cmd/config"
	"go-grpc-unary/cmd/services"
	productPb "go-grpc-unary/pb/product"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	netListen, err := net.Listen("tcp", port)

	if err != nil {
		panic(err)
	}

	db := config.ConnectDatabase();

	

	//gRPC SERVER
	grpcServer := grpc.NewServer()

	//PRODUCT SERVICE
	productService := services.ProductService{DB: db}
	productPb.RegisterProductServiceServer(grpcServer, &productService)

	fmt.Println("Server is now started at port: ", netListen.Addr())

	if err := grpcServer.Serve((netListen)); err != nil {
		panic(err)
	}

	
}