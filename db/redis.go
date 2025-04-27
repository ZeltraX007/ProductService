package db

import (
	"ProductService/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var ctx = context.Background()

type Redis struct {
	Con *redis.Client
}

func NewRedisConnector(conn *redis.Client) CacheInterface {
	return &Redis{
		Con: conn,
	}
}

// redis function to get product details by id
// GetProductByID retrieves and unmarshals the product from Redis by ID
func (r *Redis) GetProductByID(id string) (*models.Product, error) {
	log.Println("Entering GetProductByID Cache")
	result, err := r.Con.Get(ctx, id).Bytes()
	if err == redis.Nil {
		// Key does not exist
		return nil, nil
	} else if err != nil {
		// Redis error
		return nil, err
	}

	var product models.Product
	if err := json.Unmarshal([]byte(result), &product); err != nil {
		// Failed to unmarshal JSON
		return nil, errors.New("failed to unmarshal product from redis")
	}
	log.Println("Product Found in Cache")
	log.Println("Exiting GetProductByID Cache")
	return &product, nil
}

// SetProductByID stores the product in Redis with a TTL
func (r *Redis) SetProductByID(id string, product *models.Product, ttl time.Duration) error {
	log.Println("Entering SetProductByID Cache")

	// Marshal the product struct to JSON
	productJSON, err := json.Marshal(product)
	if err != nil {
		log.Println("Failed to marshal product:", err)
		return errors.New("failed to marshal product for redis")
	}

	// Store in Redis
	err = r.Con.Set(ctx, id, productJSON, ttl).Err()
	if err != nil {
		log.Println("Failed to store product in Redis:", err)
		return err
	}

	log.Println("Product Stored Successfully in Cache")
	log.Println("Exiting SetProductByID Cache")
	return nil
}

func (r *Redis) DeleteProductFromCache(id string) error {
	log.Println("Entering DeleteProductFromCache Cache")
	log.Println("Deleting product from cache:", id)
	err := r.Con.Del(ctx, id).Err()
	if err != nil && err != redis.Nil {
		return err
	}
	log.Println("Exiting DeleteProductFromCache Cache")
	return nil
}
