package db

import (
	"ProductService/models"
	"database/sql"
	"errors"
	"log"
)

type PGConnector struct {
	Conn *sql.DB
}

func NewPGConnector(conn *sql.DB) DBOperations {
	return &PGConnector{Conn: conn}
}

// Should contain all the implemented functions

func (d *PGConnector) Try() {
	log.Println("Successful interface")
}

func (d *PGConnector) GetProductByID(id int) (*models.Product, error) {
	log.Println("Entering GetProductByID DB Function")
	var product models.Product

	query := "SELECT id, name, price FROM products WHERE id = $1"
	err := d.Conn.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No product found with given ID
			return nil, nil
		}
		// Some other database error
		return nil, err
	}
	log.Println("Product Found in DB")
	log.Println("Exiting GetProductByID DB Function")
	return &product, nil
}

func (d *PGConnector) GetAllProducts(offset int, pageSize int) ([]*models.Product, error) {
	log.Println("Entering GetAllProducts DB Function")
	query := "SELECT id, name, price FROM products OFFSET $1 LIMIT $2"
	rows, err := d.Conn.Query(query, offset, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	// check for row iteration error
	if err = rows.Err(); err != nil {
		return nil, err
	}

	log.Println("Exiting GetAllProducts DB Function")
	return products, nil
}

func (d *PGConnector) CreateProduct(product *models.CreateProductRequest) (int, error) {
	log.Println("Entering CreateProduct DB Function")
	query := "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id"
	var id int
	err := d.Conn.QueryRow(query, product.Name, product.Price).Scan(&id)
	if err != nil {
		return 0, err
	}
	log.Println("Exiting CreateProduct DB Function")
	return id, nil
}

func (d *PGConnector) UpdateProduct(product *models.Product) error {
	log.Println("Entering UpdateProduct DB Function")
	query := "UPDATE products SET name = $1, price = $2 WHERE id = $3"
	result, err := d.Conn.Exec(query, product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // standard error if product not found
	}

	log.Println("Exiting UpdateProduct DB Function")
	return nil
}

func (d *PGConnector) DeleteProduct(id int) error {
	log.Println("Entering DeleteProduct DB Function")
	query := "DELETE FROM products WHERE id = $1"
	result, err := d.Conn.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Product not found
	}

	log.Println("Exiting DeleteProduct DB Function")
	return nil
}

func (d *PGConnector) GetProductCount() (int, error) {
	log.Println("Entering GetProductCount DB Function")
	query := "SELECT COUNT(*) FROM products"
	var count int
	err := d.Conn.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	log.Println("Entering GetProductCount DB Function")
	return count, nil
}
