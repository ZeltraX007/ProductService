package services_test

import (
	"ProductService/db/mocks"
	"ProductService/models"
	"ProductService/services"
	enum "ProductService/utils/enums"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduct_Decode(t *testing.T) {
	service := services.NewCreateProduct(nil, nil)

	validPayload := `{"name": "Test Product", "price": 99.99}`
	invalidPayload := `{"name": "Test Product", "price":}`

	// Valid Decode
	v, err := service.Decode([]byte(validPayload))
	assert.NoError(t, err)
	assert.NotNil(t, v)

	// Invalid Decode
	v, err = service.Decode([]byte(invalidPayload))
	assert.Error(t, err)
	assert.Nil(t, v)
}

func TestCreateProduct_Validate(t *testing.T) {
	service := services.NewCreateProduct(nil, nil)

	// Valid input
	validInput := &models.CreateProductRequest{
		Name:  "Valid Product",
		Price: 100.0,
	}
	err := service.Validate(validInput)
	assert.NoError(t, err)

	// Invalid input - Missing Name
	invalidInput := &models.CreateProductRequest{
		Name:  "",
		Price: 100.0,
	}
	err = service.Validate(invalidInput)
	assert.Error(t, err)

	// Invalid input - Price <= 0
	invalidInput2 := &models.CreateProductRequest{
		Name:  "Invalid Product",
		Price: 0.0,
	}
	err = service.Validate(invalidInput2)
	assert.Error(t, err)
}

func TestCreateProduct_ProcessMsg_Success(t *testing.T) {
	mockDB := new(mocks.MockDBOperations)
	service := services.NewCreateProduct(nil, mockDB)

	mockRequest := &models.CreateProductRequest{
		Name:  "New Product",
		Price: 50.5,
	}

	mockDB.On("CreateProduct", mock.Anything).Return(1, nil)

	resp, err := service.ProcessMsg(mockRequest, nil)

	assert.NoError(t, err)

	result, ok := resp.(models.Result)
	assert.True(t, ok)
	assert.Equal(t, enum.SuccessCode, result.ResponseCode)
	assert.Equal(t, "Product created successfully", result.ResponseDescription)

	mockDB.AssertExpectations(t)
}

func TestCreateProduct_ProcessMsg_DBError(t *testing.T) {
	mockDB := new(mocks.MockDBOperations)
	service := services.NewCreateProduct(nil, mockDB)

	mockRequest := &models.CreateProductRequest{
		Name:  "New Product",
		Price: 50.5,
	}

	mockDB.On("CreateProduct", mock.Anything).Return(0, errors.New("db error"))

	resp, err := service.ProcessMsg(mockRequest, nil)

	assert.NoError(t, err)

	result, ok := resp.(models.Result)
	assert.True(t, ok)
	assert.Equal(t, enum.FailureCode500, result.ResponseCode)
	assert.Equal(t, "Database Error", result.ResponseDescription)

	mockDB.AssertExpectations(t)
}

func TestCreateProduct_Encode(t *testing.T) {
	service := services.NewCreateProduct(nil, nil)

	result := models.Result{
		ResponseCode:        enum.SuccessCode,
		ResponseStatus:      enum.SuccessMessage,
		ResponseDescription: "Product created successfully",
	}

	data, statusCode, err := service.Encode(result)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.JSONEq(t, `{"response_code":"200","response_status":"Success","response_description":"Product created successfully","response_body":null}`, string(data))
}

func TestCreateProduct_Encode_TypeAssertionFail(t *testing.T) {
	service := services.NewCreateProduct(nil, nil)

	// Pass invalid type
	data, statusCode, err := service.Encode("invalid type")

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
	assert.Nil(t, data)
}
