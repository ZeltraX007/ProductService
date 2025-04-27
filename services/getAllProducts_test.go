package services_test

import (
	"ProductService/db/mocks"
	"ProductService/models"
	"ProductService/services"
	"ProductService/utils/enums"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestGetAllProd_ProcessMsg_DBErrorOnCount(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewGetAllProd(mockCache, mockDB)

	mockDB.On("GetProductCount").Return(0, errors.New("db error"))

	req := httptest.NewRequest("GET", "/products", nil)

	resp, err := service.ProcessMsg(nil, req)

	result := resp.(models.Result)
	assert.NoError(t, err)
	assert.Equal(t, enums.FailureCode500, result.ResponseCode)
	assert.Equal(t, enums.FailureMessage500, result.ResponseStatus)

	mockDB.AssertExpectations(t)
}

func TestGetAllProd_ProcessMsg_NoProductsFound(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewGetAllProd(mockCache, mockDB)

	mockDB.On("GetProductCount").Return(0, nil)

	req := httptest.NewRequest("GET", "/products", nil)

	resp, err := service.ProcessMsg(nil, req)

	result := resp.(models.PaginatedResponse)
	assert.NoError(t, err)
	assert.Equal(t, enums.SuccessCode, result.ResponseCode)
	assert.Equal(t, "No Products Found", result.ResponseDescription)

	mockDB.AssertExpectations(t)
}

func TestGetAllProd_ProcessMsg_ErrorInPaginationFunction(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewGetAllProd(mockCache, mockDB)

	mockDB.On("GetProductCount").Return(10, nil)

	req := httptest.NewRequest("GET", "/products?page=invalid", nil)
	q := req.URL.Query()
	q.Add("page", "invalid")
	req.URL.RawQuery = q.Encode()

	resp, err := service.ProcessMsg(nil, req)

	result := resp.(models.PaginatedResponse)
	assert.NoError(t, err)
	assert.Equal(t, enums.FailureCode500, result.ResponseCode)
	assert.Contains(t, result.ResponseDescription, "Error in PagenationFunction")

	mockDB.AssertExpectations(t)
}

func TestGetAllProd_ProcessMsg_DBErrorOnFetchProducts(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewGetAllProd(mockCache, mockDB)

	mockDB.On("GetProductCount").Return(10, nil)
	mockDB.On("GetAllProducts", 0, 10).Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/products?page=1&page_size=10", nil)
	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("page_size", "10")
	req.URL.RawQuery = q.Encode()

	resp, err := service.ProcessMsg(nil, req)

	result := resp.(models.PaginatedResponse)
	assert.NoError(t, err)
	assert.Equal(t, enums.FailureCode500, result.ResponseCode)
	assert.Equal(t, enums.FailureMessage500, result.ResponseStatus)

	mockDB.AssertExpectations(t)
}
