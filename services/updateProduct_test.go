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
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateProduct_ProcessMsg_InvalidProductId(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewUpdateProduct(mockCache, mockDB)

	req := httptest.NewRequest("PUT", "/products/invalid", strings.NewReader(""))
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

	resp, err := service.ProcessMsg(nil, req)

	result := resp.(models.Result)
	assert.NoError(t, err)
	assert.Equal(t, enums.FailureCode400, result.ResponseCode)
	assert.Equal(t, "Invalid product ID", result.ResponseDescription)

	mockDB.AssertExpectations(t)
}

func TestUpdateProduct_ProcessMsg_ProductNotFound(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewUpdateProduct(mockCache, mockDB)

	mockDB.On("UpdateProduct", mock.Anything).Return(sql.ErrNoRows)

	productReq := &models.UpdateProductRequest{
		Name:  "Updated Product",
		Price: 100,
	}

	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(""))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(productReq, req)

	result := resp.(models.Result)
	assert.NoError(t, err)
	assert.Equal(t, enums.FailureCode404, result.ResponseCode)
	assert.Equal(t, "Product not found", result.ResponseDescription)

	mockDB.AssertExpectations(t)
}

func TestUpdateProduct_ProcessMsg_DBError(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewUpdateProduct(mockCache, mockDB)

	mockDB.On("UpdateProduct", mock.Anything).Return(errors.New("db error"))

	productReq := &models.UpdateProductRequest{
		Name:  "Updated Product",
		Price: 100,
	}

	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(""))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(productReq, req)

	result := resp.(models.Result)
	assert.NoError(t, err)
	assert.Equal(t, enums.FailureCode500, result.ResponseCode)
	assert.Equal(t, "Database Error", result.ResponseDescription)

	mockDB.AssertExpectations(t)
}

func TestUpdateProduct_ProcessMsg_Success(t *testing.T) {
	mockCache := new(mocks.MockCacheInterface)
	mockDB := new(mocks.MockDBOperations)
	service := services.NewUpdateProduct(mockCache, mockDB)

	mockDB.On("UpdateProduct", mock.Anything).Return(nil)

	productReq := &models.UpdateProductRequest{
		Name:  "Updated Product",
		Price: 100,
	}

	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(""))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	resp, err := service.ProcessMsg(productReq, req)

	result := resp.(models.Result)
	assert.NoError(t, err)
	assert.Equal(t, enums.SuccessCode, result.ResponseCode)
	assert.Equal(t, "Product updated successfully", result.ResponseDescription)

	mockDB.AssertExpectations(t)
}
