package services

import (
	"ProductService/db"
	"ProductService/models"
	enum "ProductService/utils/enums"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type DeleteProd struct {
	RedisConnector db.CacheInterface
	PGDBConnector  db.DBOperations
}

func NewDeleteProd(redis db.CacheInterface, pgdb db.DBOperations) *DeleteProd {
	return &DeleteProd{
		RedisConnector: redis,
		PGDBConnector:  pgdb,
	}
}

func (b *DeleteProd) Decode(data []byte) (interface{}, error) {
	log.Printf("Entered DeleteProd Decode")
	log.Printf("Exit DeleteProd Decode")
	return nil, nil
}

func (b *DeleteProd) Validate(v interface{}) error {
	log.Printf("Entered DeleteProd Validate")
	log.Printf("Exit DeleteProd Validate")
	return nil
}

func (b *DeleteProd) ProcessMsg(v interface{}, r *http.Request) (interface{}, error) {
	log.Printf("Entered DeleteProd ProcessMsg")

	vars := mux.Vars(r)
	productIdStr := vars["id"]
	log.Println("Product ID:", productIdStr)
	// Validate product ID
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		msg := models.Result{
			ResponseCode:        enum.FailureCode400,
			ResponseStatus:      enum.FailureMessage400,
			ResponseDescription: enum.FailureMessage400,
			ResponseBody:        nil,
		}
		return msg, nil
	}

	// Delete product from database
	err = b.PGDBConnector.DeleteProduct(productId)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := models.Result{
				ResponseCode:        enum.FailureCode404,
				ResponseStatus:      enum.FailureMessage404,
				ResponseDescription: "Product not found",
				ResponseBody:        nil,
			}
			return msg, nil
		}

		msg := models.Result{
			ResponseCode:        enum.FailureCode500,
			ResponseStatus:      enum.FailureMessage500,
			ResponseDescription: "Database Error",
			ResponseBody:        nil,
		}
		return msg, nil
	}

	// Also, optionally delete from Redis cache
	err = b.RedisConnector.DeleteProductFromCache(productIdStr)
	if err != nil {
		log.Printf("Failed to delete product from cache: %v", err)
	}

	msg := models.Result{
		ResponseCode:        enum.SuccessCode,
		ResponseStatus:      enum.SuccessMessage,
		ResponseDescription: "Product deleted successfully",
		ResponseBody:        nil,
	}
	return msg, nil
}

func (b *DeleteProd) Encode(v interface{}) ([]byte, int, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic occurred: %v", r)
		}
	}()
	log.Printf("Entered CreateProduct Encode")

	format, ok := v.(models.Result)
	if !ok {
		log.Printf("Type assertion failed: expected models.Result but got %T", v)
		return nil, http.StatusInternalServerError, fmt.Errorf("type assertion failed: expected models.Result but got %T", v)
	}

	data, err := json.Marshal(&format)
	if err != nil {
		log.Println("Error in Marshal", err)
		return nil, http.StatusInternalServerError, err
	}

	// Decide HTTP status code based on ResponseCode
	statusCode := http.StatusOK // default 200
	switch format.ResponseCode {
	case "400":
		statusCode = http.StatusBadRequest
	case "404":
		statusCode = http.StatusNotFound
	case "500":
		statusCode = http.StatusInternalServerError
	}

	log.Printf("Exit CreateProduct Encode")
	return data, statusCode, nil
}
