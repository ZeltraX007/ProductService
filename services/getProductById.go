package services

import (
	"ProductService/db"
	"ProductService/models"
	enum "ProductService/utils/enums"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type GetProdById struct {
	RedisConnector db.CacheInterface
	PGDBConnector  db.DBOperations
}

func NewGetProdById(redis db.CacheInterface, pgdb db.DBOperations) *GetProdById {
	return &GetProdById{
		RedisConnector: redis,
		PGDBConnector:  pgdb,
	}
}

func (b *GetProdById) Decode(data []byte) (interface{}, error) {
	log.Printf("Entered GetProdById Decode")
	log.Printf("Exit GetProdById Decode")
	return nil, nil
}

func (b *GetProdById) Validate(v interface{}) error {
	log.Printf("Entered GetProdById Validate")
	log.Printf("Exit GetProdById Validate")
	return nil
}

func (b *GetProdById) ProcessMsg(v interface{}, r *http.Request) (interface{}, error) {
	log.Printf("Entered GetProdById ProcessMsg")
	var product *models.Product
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
	log.Printf("done 1")
	product, err = b.RedisConnector.GetProductByID(productIdStr)
	if err != nil {
		msg := models.Result{
			ResponseCode:        enum.FailureCode500,
			ResponseStatus:      enum.FailureMessage500,
			ResponseDescription: enum.FailureMessage500,
			ResponseBody:        nil,
		}
		return msg, nil
	}
	log.Printf("done 2")

	if product == nil {
		log.Println("Product Not Found in Cache")
		product, err = b.PGDBConnector.GetProductByID(productId)
		if err != nil {
			msg := models.Result{
				ResponseCode:        enum.FailureCode500,
				ResponseStatus:      enum.FailureMessage500,
				ResponseDescription: "Database Error",
				ResponseBody:        nil,
			}
			return msg, nil
		}

		if product != nil {
			err = b.RedisConnector.SetProductByID(productIdStr, product, time.Minute)
			if err != nil {
				return nil, err
			}
		}
	}

	if product == nil {
		log.Println("Product Not Found in DB")
		msg := models.Result{
			ResponseCode:        enum.FailureCode404,
			ResponseStatus:      enum.FailureMessage404,
			ResponseDescription: "Product Not Found",
			ResponseBody:        nil,
		}
		return msg, nil
	}

	msg := models.Result{
		ResponseCode:        enum.SuccessCode,
		ResponseStatus:      enum.SuccessMessage,
		ResponseDescription: "Product fetched successfully",
		ResponseBody:        product,
	}
	log.Printf("Exiting GetProdById ProcessMsg")
	return msg, nil
}

func (b *GetProdById) Encode(v interface{}) ([]byte, int, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic occurred: %v", r)
		}
	}()
	log.Printf("Entered GetProdById Encode")

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

	log.Printf("Exit GetProdById Encode")
	return data, statusCode, nil
}
