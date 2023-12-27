package services

import (
	"context"
	productPb "go-grpc-unary/pb/product"
	pagingPb "go-grpc-unary/pb/pagination"
)
type ProductService struct {
	productPb.UnimplementedProductServiceServer
}

func (p *ProductService) GetProducts(context.Context, *productPb.Empty) (*productPb.Products, error) {
	products := &productPb.Products{
		Pagination: &pagingPb.Pagination{
			Total: 10,
			PerPage: 4,
			CurrentPage: 1,
			LastPage: 2,
		},
		Data: []*productPb.Product{
			{
				Id: 1,
				Name: "Nike Tees",
				Price: 10000,
				Stock: 10,
				Category: &productPb.Category{
					Id: 1,
					Name: "shirt",
				},
			},
			{
				Id: 2,
				Name: "Adidas Tees",
				Price: 120000,
				Stock: 3,
				Category: &productPb.Category{
					Id: 1,
					Name: "shirt",
				},
			},
		},
	}

	return products, nil
}