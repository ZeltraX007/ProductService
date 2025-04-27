package services

import (
	"ProductService/db"
	"ProductService/models"
	enum "ProductService/utils/enums"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type CreateProduct struct {
	RedisConnector db.CacheInterface
	PGDBConnector  db.DBOperations
}

func NewCreateProduct(redis db.CacheInterface, pgdb db.DBOperations) *CreateProduct {
	return &CreateProduct{
		RedisConnector: redis,
		PGDBConnector:  pgdb,
	}
}

func (b *CreateProduct) Decode(data []byte) (interface{}, error) {
	log.Printf("Entered CreateProduct Decode")
	var format *models.CreateProductRequest
	err := json.Unmarshal(data, &format)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("Exit CreateProduct Decode")
	return format, nil
}

func (b *CreateProduct) Validate(v interface{}) error {
	log.Printf("Entered CreateProduct Validate")
	format := v.(*models.CreateProductRequest)
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
	log.Printf("Exit CreateProduct Validate")
	return nil
}

func (b *CreateProduct) ProcessMsg(v interface{}, r *http.Request) (interface{}, error) {
	log.Println("Entered CreateProduct ProcessMsg")

	product := v.(*models.CreateProductRequest)

	// Create product in database
	_, err := b.PGDBConnector.CreateProduct(product)
	if err != nil {
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
		ResponseDescription: "Product created successfully",
		ResponseBody:        nil,
	}
	log.Println("Exiting CreateProduct ProcessMsg")
	return msg, nil
}

func (b *CreateProduct) Encode(v interface{}) ([]byte, int, error) {
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
