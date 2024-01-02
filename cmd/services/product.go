package services

import (
	"context"
	"go-grpc-unary/cmd/helpers"
	pagingPb "go-grpc-unary/pb/pagination"
	productPb "go-grpc-unary/pb/product"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)
type ProductService struct {
	productPb.UnimplementedProductServiceServer
	DB *gorm.DB
}


//GETTING ALL PRODUCTS
func (p *ProductService) GetProducts(ctx context.Context, pageParam *productPb.Page) (*productPb.Products, error) {
	//DYNAMIC PAGINATION FUNCTION
	var page int64 = 1
	if pageParam.GetPage() != 0 {
		page = pageParam.GetPage()
	}

	var pagination pagingPb.Pagination
	var products []*productPb.Product

	//SQL QUERY FOR GETTING ALL DATA
	sql := p.DB.Table("products as p").
	Joins("INNER JOIN categories AS c on c.id = p.category_id").
	Select("p.id", "p.name", "p.price", "p.stock", "c.id", "c.name")

	//SETTING OFFSET AND LIMIT FOR EACH PAGE -> USED IN PAGINATION
	offset, limit := helpers.Pagination(sql, page, &pagination)

	//APPLYING OFFSET AND LIMIT FOR QUERY 
	rows, err := sql.Offset(int(offset)).Limit(int(limit)).Rows()

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	//INSERTING EACH ROW TO DATABASE
	for rows.Next() {
		var product productPb.Product
		var category productPb.Category

		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Stock, &category.Id, &category.Name); err != nil {
			log.Fatal("error getting data")
		}

		product.Category = &category
		products = append(products, &product)	
	}

	//SETTING UP RESPONSES
	response := &productPb.Products{
		Pagination: &pagination,
		Data: products,
	}

	return response, nil

}


//GETTING SPECIFIC PRODUCT
func (p *ProductService) GetProduct(ctx context.Context, id *productPb.Id) (*productPb.Product, error) {
	row := p.DB.Table("products as p").
	Joins("INNER JOIN categories AS c on c.id = p.category_id").
	Select("p.id", "p.name", "p.price", "p.stock", "c.id", "c.name").
	Where("p.id = ?", id.GetId()).
	Row()

	var product productPb.Product
	var category productPb.Category

	//READING THE DATA
	if err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Stock, &category.Id, &category.Name); err != nil {
		log.Fatal("error getting data")
	}

	product.Category = &category
	
	return &product, nil
}

func(p *ProductService) CreateProduct(ctx context.Context, productData *productPb.Product) (*productPb.Id, error) {
	var Response productPb.Id

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		//GETTING CATEGORY
		category := productPb.Category{
			Id: 0,
			Name: productData.GetCategory().GetName(),
		}

		if err := tx.Table("categories").
		Where("LCASE(name) = ?", category.GetName()).
		FirstOrCreate(&category).Error; err != nil {
			return err
		}

		//NEW PRODUCT STRUCT
		product := struct {
			Id uint64
			Name string
			Price uint64
			Stock uint64
			Category_id uint64
		} {
			Id: productData.GetId(),
			Name: productData.GetName(),
			Price: productData.GetPrice(),
			Stock: productData.GetStock(),
			Category_id: category.GetId(),
		}

		if err := tx.Table("products").Create(&product).Error; err != nil {
			return err
		}

		Response.Id = product.Id
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &Response, nil
}


//UPDATING SPECIFIC PRODUCT
func (p *ProductService) UpdateProduct(ctx context.Context, productData *productPb.Product) (*productPb.Status, error) {
	var Response productPb.Status

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		category := productPb.Category{
			Id: 0,
			Name: productData.GetCategory().GetName(),
		}

		if err := tx.Table("categories").
		Where("LCASE(name) = ?", category.GetName()).
		FirstOrCreate(&category).Error; err != nil {
			return err
		}

		product := struct {
			Id uint64
			Name string
			Price uint64
			Stock uint64
			Category_id uint64
		} {
			Id: productData.GetId(),
			Name: productData.GetName(),
			Price: productData.GetPrice(),
			Stock: productData.GetStock(),
			Category_id: category.GetId(),
		}

		if err := tx.Table("products").Where("id = ?", product.Id).Updates(&product).Error; err != nil {
			return err
		}

		Response.Status = 1
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &Response, nil
}

//DELETING SPECIFIC DATA
func (p *ProductService) DeleteProduct(ctx context.Context, id *productPb.Id) (*productPb.Status, error) {
	var response productPb.Status

	if err := p.DB.Table("products").Where("id = ?", id.GetId()).Delete(nil).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response.Status = 1
	return &response, nil
}