package services

import (
	"ProductService/db"
	"ProductService/models"
	enum "ProductService/utils/enums"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type UpdateProduct struct {
	RedisConnector db.CacheInterface
	PGDBConnector  db.DBOperations
}

func NewUpdateProduct(redis db.CacheInterface, pgdb db.DBOperations) *UpdateProduct {
	return &UpdateProduct{
		RedisConnector: redis,
		PGDBConnector:  pgdb,
	}
}

func (b *UpdateProduct) Decode(data []byte) (interface{}, error) {
	log.Printf("Entered UpdateProduct Decode")
	var format *models.UpdateProductRequest
	err := json.Unmarshal(data, &format)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("Exit UpdateProduct Decode")
	return format, nil
}

func (b *UpdateProduct) Validate(v interface{}) error {
	log.Printf("Entered UpdateProduct Validate")
	format := v.(*models.UpdateProductRequest)
	var validate = validator.New()
	e := validate.Struct(v)
	if e != nil {
		log.Println(e)
		return e
	}

	if format.Name == "" || format.Price <= 0 {
		err := errors.New("mandatory fields are missing in request")
		return err
	}
	log.Printf("Exit UpdateProduct Validate")
	return nil
}

func (b *UpdateProduct) ProcessMsg(v interface{}, r *http.Request) (interface{}, error) {
	log.Println("Entered UpdateProduct ProcessMsg")
	// Extract product ID from URL query
	vars := mux.Vars(r)
	productIdStr := vars["id"]
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		msg := models.Result{
			ResponseCode:        enum.FailureCode400,
			ResponseStatus:      enum.FailureMessage400,
			ResponseDescription: "Invalid product ID",
			ResponseBody:        nil,
		}
		return msg, nil
	}

	product := v.(*models.UpdateProductRequest)

	updatedProduct := models.Product{
		ID:    productId,
		Name:  product.Name,
		Price: product.Price,
	}

	err = b.PGDBConnector.UpdateProduct(&updatedProduct)
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

	msg := models.Result{
		ResponseCode:        enum.SuccessCode,
		ResponseStatus:      enum.SuccessMessage,
		ResponseDescription: "Product updated successfully",
		ResponseBody:        nil,
	}
	return msg, nil
}

func (b *UpdateProduct) Encode(v interface{}) ([]byte, int, error) {
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
