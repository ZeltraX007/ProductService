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
	"math"
	"net/http"
	"strconv"
)

type GetAllProd struct {
	RedisConnector db.CacheInterface
	PGDBConnector  db.DBOperations
}

func NewGetAllProd(redis db.CacheInterface, pgdb db.DBOperations) *GetAllProd {
	return &GetAllProd{
		RedisConnector: redis,
		PGDBConnector:  pgdb,
	}
}

func (b *GetAllProd) Decode(data []byte) (interface{}, error) {
	log.Printf("Entered GetAllProd Decode")
	log.Printf("Exit GetAllProd Decode")
	return nil, nil
}

func (b *GetAllProd) Validate(v interface{}) error {
	log.Printf("Entered GetAllProd Validate")
	log.Printf("Exit GetAllProd Validate")
	return nil
}

func (b *GetAllProd) ProcessMsg(v interface{}, r *http.Request) (interface{}, error) {
	log.Println("Entered GetAllProd ProcessMsg")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	var emptyResponse models.PaginationProductResponse

	count, err := b.PGDBConnector.GetProductCount()
	if err != nil {
		msg := models.Result{
			ResponseCode:        enum.FailureCode500,
			ResponseStatus:      enum.FailureMessage500,
			ResponseDescription: "Database Error",
			ResponseBody:        nil,
		}
		return msg, nil
	}

	if count == 0 {
		msg := models.PaginatedResponse{
			ResponseCode:        enum.SuccessCode,
			ResponseStatus:      enum.SuccessMessage,
			ResponseDescription: "No Products Found",
			ResponseBody:        emptyResponse,
		}
		return msg, nil
	}

	pageBody, e := PagenationFunction(pageStr, pageSizeStr, count)
	if e != nil {
		log.Println("Error in PagenationFunction: ", e)
		msg := models.PaginatedResponse{
			ResponseCode:        enum.FailureCode500,
			ResponseStatus:      enum.FailureMessage500,
			ResponseDescription: "Error in PagenationFunction",
			ResponseBody:        emptyResponse,
		}
		return msg, nil
	}

	pageBodyResp := pageBody.(models.PaginationProductResponse)

	products, err := b.PGDBConnector.GetAllProducts(pageBodyResp.Offset, pageBodyResp.PageSize)
	if err != nil {
		msg := models.PaginatedResponse{
			ResponseCode:        enum.FailureCode500,
			ResponseStatus:      enum.FailureMessage500,
			ResponseDescription: "Database Error",
			ResponseBody:        emptyResponse,
		}
		return msg, nil
	}

	response := models.PaginationProductResponse{
		PageNo:     pageBodyResp.PageNo,
		PageSize:   pageBodyResp.PageSize,
		TotalCount: pageBodyResp.TotalCount,
		TotalPages: pageBodyResp.TotalPages,
		Offset:     pageBodyResp.Offset,
		Products:   products,
	}

	msg := models.PaginatedResponse{
		ResponseCode:        enum.SuccessCode,
		ResponseStatus:      enum.SuccessMessage,
		ResponseDescription: "Products fetched successfully",
		ResponseBody:        response,
	}
	return msg, nil
}

func (b *GetAllProd) Encode(v interface{}) ([]byte, int, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic occurred: %v", r)
		}
	}()
	log.Printf("Entered GetAllProd Encode")

	format, ok := v.(models.PaginatedResponse)
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

	log.Printf("Exit GetAllProd Encode")
	return data, statusCode, nil
}

func PagenationFunction(received_page_str interface{}, received_page_size_str interface{}, received_total_elements_str interface{}) (interface{}, error) {
	log.Println("Entered PagenationFunction")

	var page, page_size int

	page_str := received_page_str.(string)
	page_size_str := received_page_size_str.(string)
	total_elements := received_total_elements_str.(int)

	if page_str == "" {
		page = 1
		log.Println("Setting default page to 1")
	} else {
		validate := validator.New()
		log.Println("Page string is :", page_str)
		res := validate.Var(page_str, "numeric")
		if res != nil {
			log.Println("Error in validating page no")
			err := errors.New("Error in validating pageno. Page number not numeric")
			return nil, err
		} else {
			var e error
			page, e = strconv.Atoi(page_str)
			if e != nil {
				log.Println("Error in converting page number to integer :", page)
				err := errors.New("Error in converting page number to integer")
				return nil, err
			}
			if page == 0 {
				page = 1
			}
		}
	}

	if page_size_str == "" {
		page_size = 10
		log.Println("Setting default page size to 10")
	} else {
		validate := validator.New()
		log.Println("Page size string is :", page_size_str)
		res := validate.Var(page_size_str, "numeric")
		if res != nil {
			log.Println("Error in validating page size")
			err := errors.New("Error in validating page size. Page size not numeric")
			return nil, err
		} else {
			var e error
			page_size, e = strconv.Atoi(page_size_str)
			log.Println("page size is", page_size)
			if e != nil {
				log.Println("Error in converting page size to integer :", page)
				err := errors.New("Error in converting page size to integer")
				return nil, err
			}
			if page_size == 0 {
				page_size = 10
			}
		}
	}

	offset := (page * page_size) - page_size

	//calculating the total number of pages
	var total_pages float64
	if page_size != 0 {
		//division by 0
		var total_count = float64(total_elements)
		var total_size = float64(page_size)
		total_pages = total_count / total_size
	} else {
		total_pages = float64(total_elements)
	}
	total_pages = math.Ceil(total_pages)
	var total_pages_req = int(total_pages)

	resp := models.PaginationProductResponse{
		PageNo:     page,
		PageSize:   page_size,
		TotalPages: total_pages_req,
		TotalCount: total_elements,
		Offset:     offset,
	}

	log.Println("Exited PagenationFunction")
	return resp, nil
}
