package services

import (
	"context"
	pagingPb "go-grpc-unary/pb/pagination"
	productPb "go-grpc-unary/pb/product"
	"log"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)
type ProductService struct {
	productPb.UnimplementedProductServiceServer
	DB *gorm.DB
}

func (p *ProductService) GetProducts(context.Context, *productPb.Empty) (*productPb.Products, error) {
	
	//DYNAMIC PAGINATION FUNCTION
	var page int64 = 1
	var total int64
	var limit int64 = 1
	var offset int64

	var pagination pagingPb.Pagination


	var products []*productPb.Product

	sql := p.DB.Table("products as p").
	Joins("INNER JOIN categories AS c on c.id = p.category_id").
	Select("p.id", "p.name", "p.price", "p.stock", "c.id", "c.name")

	sql.Count(&total)
	if page == 1 {
		offset = 0
	} else {
		offset = (page - 1) * limit
	}

	pagination.Total = uint32(total)
	pagination.PerPage = uint32(limit)
	pagination.CurrentPage = uint32(page)
	pagination.LastPage = uint32(math.Ceil(float64(total))/ float64(limit))

	rows, err := sql.Offset(int(offset)).Limit(int(limit)).Rows()


	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var product productPb.Product
		var category productPb.Category

		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Stock, &category.Id, &category.Name); err != nil {
			log.Fatal("error getting data")
		}

		product.Category = &category
		products = append(products, &product)	
	}

	response := &productPb.Products{
		Pagination: &pagination,
		Data: products,
	}

	return response, nil

}