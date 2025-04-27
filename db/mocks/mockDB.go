package mocks

import (
	"ProductService/models"
	"github.com/stretchr/testify/mock"
)

type MockDBOperations struct {
	mock.Mock
}

func (m *MockDBOperations) Try() {
	m.Called()
}

func (m *MockDBOperations) GetProductByID(id int) (*models.Product, error) {
	args := m.Called(id)
	product, ok := args.Get(0).(*models.Product)
	if !ok {
		return nil, args.Error(1)
	}
	return product, args.Error(1)
}

func (m *MockDBOperations) GetAllProducts(offset int, pageSize int) ([]*models.Product, error) {
	args := m.Called(offset, pageSize)
	products, ok := args.Get(0).([]*models.Product)
	if !ok {
		return nil, args.Error(1)
	}
	return products, args.Error(1)
}

func (m *MockDBOperations) CreateProduct(product *models.CreateProductRequest) (int, error) {
	args := m.Called(product)
	return args.Int(0), args.Error(1)
}

func (m *MockDBOperations) UpdateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockDBOperations) DeleteProduct(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDBOperations) GetProductCount() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}
