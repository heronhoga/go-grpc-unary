package services

import (
	productPb "go-grpc-unary/pb/product"
)
type ProductService struct {
	productPb.UnimplementedProductServiceServer
}