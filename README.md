# 🛍️ ProductService

A lightweight Go microservice to **manage products** — including CRUD operations on product details — with **PostgreSQL** database and **Redis** cache integration.

---

## 📚 Table of Contents

- [Features](#-features)
- [Project Structure](#-project-structure)
- [Setup Instructions](#-setup-instructions)
  - [Prerequisites](#prerequisites)
- [API Overview](#-api-endpoints)
    - [Create Product](#create-product)
    - [Get Product By ID](#get-product-by-id)
    - [Get All Products](#get-all-products)
    - [Update Product](#update-product)
    - [Delete Product](#delete-product)
- [Error Handling](#-error-handling)
- [Testing](#-testing)

---

## 📦 Features

- **List Products** (with smart pagination)
- **Create Products**
- **Update Product** (by product ID)
- **Delete Product**
- **Clean and modular service structure**
- **Error handling and validation**
- **Interfaces for easy mocking** during testing

---

## 🏗️ Project Structure

```
ProductService/
├── app/              # Initializing the project
├── db/               # Database and cache connectors
    └── mocks/        # mocks for unit testing
├── models/           # Request and response models
├── services/         # Business logic (GetAllProd, UpdateProduct, etc.)
├── utils/enums/      # Common enums and constants
└── main.go           # Application start
```

---

## ⚙️ Setup instructions

###  ️Prerequisites

Before you begin, ensure you have the following installed:
- **Go (Golang)**: The Go programming language. [Install Go](https://golang.org/dl/)
- **PostgreSQL**: A relational database system. [Install PostgreSQL](https://www.postgresql.org/download/)
- **Redis**: An in-memory data structure store. [Install Redis](https://redis.io/download)
- **Go modules**: Go's dependency management system (should be installed automatically with Go).

You also need to have access to a PostgreSQL database and a Redis instance.

1. **Clone the repository**
   ```bash
   git clone https://github.com/ZeltraX007/ProductService.git
   cd ProductService
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure environment variables**
    - Create a `.env` file or set directly:
      ```
      PG_HOST=localhost
      PG_PORT=5432
      PG_USER = "postgres"
      PG_PASS = "12345"
      PG_DBNAME = "postgres"
      REDIS_HOST="localhost"
      REDIS_PORT="6379"
      REDIS_PASS=
      PORT="8000"
      HTTP_CLIENT_TIMEOUT="60"
      ```

4. **Run the application**
   ```bash
   go run main.go
   ```

---

## 🔥 API Endpoints

| Method | Endpoint                        | Description                                  |
|:-------|:--------------------------------|:---------------------------------------------|
| GET    | `/products?page=1&page_size=10` | Fetches paginated products list              |
| GET    | `/products/{id}`                | Fetches products by id                       |
| POST   | `/products`                     | Creates a product and inserts it in database |
| PUT    | `/products/{id}`          | Update an existing product                   |
| DELETE | `/products/{id}`          | Deletes an existing product                  |

---

### Create Product

```http
POST /products
```

- **Request body**: A JSON object containing `name`, `price`, etc.
- **Response**: Returns the created product with HTTP status `201 Created`.

### Get Product By ID

```http
GET /products/{id}
```

- **URL Parameter**: `id` (Product ID)
- **Response**: Returns the product details if found, or a `404 Not Found` error if not.

### Get All Products

```http
GET /products
```

- **Query Parameters**:
    - `page` (default: 1)
    - `page_size` (default: 10)
- **Response**: Returns a paginated list of products.

### Update Product

```http
PUT /products/{id}
```

- **URL Parameter**: `id` (Product ID)
- **Request body**: A JSON object containing updated `name`, `price`, etc.
- **Response**: Returns a success message upon successful update.

### Delete Product

```http
DELETE /products/{id}
```

- **URL Parameter**: `id` (Product ID)
- **Response**: Returns a success message upon deletion or a `404 Not Found` error if the product doesn't exist.

---

## ⚠️ Error Handling

The **Product Service API** follows standardized error handling practices to ensure consistent and predictable responses for clients.

### Error Codes

- **400 - Bad Request**  
  Indicates that the server could not understand the request due to invalid syntax or missing/incorrect fields.

- **404 - Not Found**  
  Returned when the requested resource (e.g., a product by ID) does not exist.

- **500 - Internal Server Error**  
  Indicates that the server encountered an unexpected condition that prevented it from fulfilling the request.

### Error Response Format

Each error response is structured as a JSON object containing:

```json
{
  "ResponseCode": "<Error Code>",
  "ResponseStatus": "<Status Message>",
  "ResponseDescription": "<Detailed Error Description>",
  "ResponseBody": null
}
```

### Example Error Responses

**400 Bad Request:**
```json
{
  "ResponseCode": "400",
  "ResponseStatus": "Failure",
  "ResponseDescription": "Invalid product ID",
  "ResponseBody": null
}
```

**404 Not Found:**
```json
{
  "ResponseCode": "404",
  "ResponseStatus": "Failure",
  "ResponseDescription": "Product not found",
  "ResponseBody": null
}
```

**500 Internal Server Error:**
```json
{
  "ResponseCode": "500",
  "ResponseStatus": "Failure",
  "ResponseDescription": "Database Error",
  "ResponseBody": null
}
```

### Notes
- Clients should handle different HTTP status codes appropriately.
- Always check the `ResponseDescription` field for a more detailed explanation of the error.

---

## 🧪 Testing

```bash
go test ./...
```

- Uses `stretchr/testify` for assertions.
- Mocks database interactions for pure unit tests.

---

## 🚀 Technologies Used

- **Go 1.20+**
- **PostgreSQL**
- **Redis**
- **Gorilla Mux** (Routing)
- **Go Playground Validator** (Validation)
- **Stretchr/testify** (Unit Testing)
- **Mockery** (Mock generation)

---

