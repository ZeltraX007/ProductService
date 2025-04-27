package services_test

import (
	"ProductService/db/mocks"
	"ProductService/models"
	"ProductService/services"
	"ProductService/utils/enums"
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestDeleteProd_ProcessMsg_InvalidProductID(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewDeleteProd(mockCache, mockDB)

	req := httptest.NewRequest("DELETE", "/products/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

	resp, err := service.ProcessMsg(nil, req)

	result := resp.(models.Result)
	assert.NoError(t, err)
	assert.Equal(t, enums.FailureCode400, result.ResponseCode)
	assert.Equal(t, enums.FailureMessage400, result.ResponseStatus)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestDeleteProd_ProcessMsg_ProductNotFoundInDB(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewDeleteProd(mockCache, mockDB)

	mockDB.On("DeleteProduct", 1).Return(sql.ErrNoRows)

	req := httptest.NewRequest("DELETE", "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(nil, req)

	result := resp.(models.Result)
	assert.NoError(t, err)
	assert.Equal(t, enums.FailureCode404, result.ResponseCode)
	assert.Equal(t, enums.FailureMessage404, result.ResponseStatus)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestDeleteProd_ProcessMsg_DBError(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewDeleteProd(mockCache, mockDB)

	mockDB.On("DeleteProduct", 1).Return(errors.New("db failure"))

	req := httptest.NewRequest("DELETE", "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(nil, req)

	result := resp.(models.Result)
	assert.NoError(t, err)
	assert.Equal(t, enums.FailureCode500, result.ResponseCode)
	assert.Equal(t, enums.FailureMessage500, result.ResponseStatus)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestDeleteProd_ProcessMsg_SuccessfulDelete(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewDeleteProd(mockCache, mockDB)

	mockDB.On("DeleteProduct", 1).Return(nil)
	mockCache.On("DeleteProductFromCache", "1").Return(nil)

	req := httptest.NewRequest("DELETE", "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(nil, req)

	result := resp.(models.Result)
	assert.NoError(t, err)
	assert.Equal(t, enums.SuccessCode, result.ResponseCode)
	assert.Equal(t, enums.SuccessMessage, result.ResponseStatus)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestDeleteProd_ProcessMsg_SuccessfulDelete_RedisError(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewDeleteProd(mockCache, mockDB)

	mockDB.On("DeleteProduct", 1).Return(nil)
	mockCache.On("DeleteProductFromCache", "1").Return(errors.New("redis delete failure"))

	req := httptest.NewRequest("DELETE", "/products/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(nil, req)

	result := resp.(models.Result)
	assert.NoError(t, err)
	assert.Equal(t, enums.SuccessCode, result.ResponseCode)
	assert.Equal(t, enums.SuccessMessage, result.ResponseStatus)

	mockCache.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}
