package mocks

import (
	"ProductService/models"
	"github.com/stretchr/testify/mock"
	"time"
)

// MockCacheInterface mocks CacheInterface
type MockCacheInterface struct {
	mock.Mock
}

func (m *MockCacheInterface) GetProductByID(id string) (*models.Product, error) {
	args := m.Called(id)
	product, ok := args.Get(0).(*models.Product)
	if !ok {
		return nil, args.Error(1)
	}
	return product, args.Error(1)
}

func (m *MockCacheInterface) SetProductByID(id string, product *models.Product, ttl time.Duration) error {
	args := m.Called(id, product, ttl)
	return args.Error(0)
}

func (m *MockCacheInterface) DeleteProductFromCache(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
