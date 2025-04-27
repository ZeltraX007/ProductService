package services_test

import (
	"ProductService/db/mocks"
	"ProductService/models"
	"ProductService/services"
	enum "ProductService/utils/enums"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetProdById_ProcessMsg_CacheHit(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewGetProdById(mockCache, mockDB)

	product := &models.Product{ID: 1, Name: "Cached Product", Price: 100}

	// Set up mocks
	mockCache.On("GetProductByID", "1").Return(product, nil)

	// Create a fake request with ID in URL
	req := httptest.NewRequest("GET", "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(nil, req)

	assert.NoError(t, err)
	result := resp.(models.Result)
	assert.Equal(t, enum.SuccessCode, result.ResponseCode)
	assert.Equal(t, product, result.ResponseBody)

	mockCache.AssertExpectations(t)
}

func TestGetProdById_ProcessMsg_CacheMiss_DBHit(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewGetProdById(mockCache, mockDB)

	product := &models.Product{ID: 1, Name: "DB Product", Price: 200}

	// Set up mocks
	mockCache.On("GetProductByID", "1").Return(nil, nil)
	mockDB.On("GetProductByID", 1).Return(product, nil)
	mockCache.On("SetProductByID", "1", product, time.Minute).Return(nil)

	req := httptest.NewRequest("GET", "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(nil, req)

	assert.NoError(t, err)
	result := resp.(models.Result)
	assert.Equal(t, enum.SuccessCode, result.ResponseCode)
	assert.Equal(t, product, result.ResponseBody)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestGetProdById_ProcessMsg_CacheMiss_DBNotFound(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewGetProdById(mockCache, mockDB)

	// Set up mocks
	mockCache.On("GetProductByID", "1").Return(nil, nil)
	mockDB.On("GetProductByID", 1).Return(nil, nil)

	req := httptest.NewRequest("GET", "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(nil, req)

	assert.NoError(t, err)
	result := resp.(models.Result)
	assert.Equal(t, enum.FailureCode404, result.ResponseCode)
	assert.Equal(t, "Product Not Found", result.ResponseDescription)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestGetProdById_ProcessMsg_InvalidProductID(t *testing.T) {
	service := services.NewGetProdById(nil, nil)

	req := httptest.NewRequest("GET", "/products/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

	resp, err := service.ProcessMsg(nil, req)

	assert.NoError(t, err)
	result := resp.(models.Result)
	assert.Equal(t, enum.FailureCode400, result.ResponseCode)
	assert.Equal(t, enum.FailureMessage400, result.ResponseDescription)
}

func TestGetProdById_ProcessMsg_CacheError(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewGetProdById(mockCache, mockDB)

	mockCache.On("GetProductByID", "1").Return(nil, errors.New("cache error"))

	req := httptest.NewRequest("GET", "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(nil, req)

	assert.NoError(t, err)
	result := resp.(models.Result)
	assert.Equal(t, enum.FailureCode500, result.ResponseCode)
	assert.Equal(t, enum.FailureMessage500, result.ResponseDescription)
}

func TestGetProdById_ProcessMsg_DBError(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewGetProdById(mockCache, mockDB)

	mockCache.On("GetProductByID", "1").Return(nil, nil)
	mockDB.On("GetProductByID", 1).Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(nil, req)

	assert.NoError(t, err)
	result := resp.(models.Result)
	assert.Equal(t, enum.FailureCode500, result.ResponseCode)
	assert.Equal(t, "Database Error", result.ResponseDescription)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestGetProdById_Encode(t *testing.T) {
	service := services.NewGetProdById(nil, nil)

	result := models.Result{
		ResponseCode:        enum.SuccessCode,
		ResponseStatus:      enum.SuccessMessage,
		ResponseDescription: "Product fetched successfully",
		ResponseBody: &models.Product{
			ID:    1,
			Name:  "Test Product",
			Price: 99.99,
		},
	}

	data, statusCode, err := service.Encode(result)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.True(t, bytes.Contains(data, []byte("Test Product")))
}

func TestGetProdById_Encode_TypeAssertionFail(t *testing.T) {
	service := services.NewGetProdById(nil, nil)

	data, statusCode, err := service.Encode("wrong type")

	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
}

func TestGetProdById_ProcessMsg_CacheMiss_DBHit_SetProductError(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewGetProdById(mockCache, mockDB)

	product := &models.Product{ID: 1, Name: "DB Product", Price: 200}

	// Set up mocks
	mockCache.On("GetProductByID", "1").Return(nil, nil)
	mockDB.On("GetProductByID", 1).Return(product, nil)
	mockCache.On("SetProductByID", "1", product, time.Minute).Return(errors.New("redis set error"))

	req := httptest.NewRequest("GET", "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(nil, req)

	// In this case, we expect an error to be returned directly
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.EqualError(t, err, "redis set error")

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}
