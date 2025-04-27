package db

import (
	"ProductService/models"
	"time"
)

type CacheInterface interface {
	GetProductByID(id string) (*models.Product, error)
	SetProductByID(id string, product *models.Product, ttl time.Duration) error
	DeleteProductFromCache(id string) error
}
