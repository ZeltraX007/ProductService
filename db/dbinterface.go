package db

import "ProductService/models"

type DBOperations interface {
	// Should contain all the postgres operations
	Try()
	GetProductByID(id int) (*models.Product, error)
	GetAllProducts(offset int, pageSize int) ([]*models.Product, error)
	CreateProduct(product *models.CreateProductRequest) (int, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id int) error
	GetProductCount() (int, error)
}
