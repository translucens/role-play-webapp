package database

import (
	"context"
	"database/sql"
)

type DevDatabaseHandler struct {
	DB *sql.DB
}

func NewDevDatabaseHandler(db *sql.DB) DevDatabaseHandler {
	return DevDatabaseHandler{DB: db}
}

func (dbh DevDatabaseHandler) InitDatabase() error {
	return nil
}

func (dbh DevDatabaseHandler) GetProduct(id int) (Product, error) {
	product := Product{
		ID: 1, Name: "product1", Price: 100, Image: "image/product1.png",
	}

	return product, nil
}

func (dbh DevDatabaseHandler) GetProducts() ([]Product, error) {
	products := []Product{
		{ID: 1, Name: "product1", Price: 100, Image: "image/product1.png"},
		{ID: 2, Name: "product2", Price: 200, Image: "image/product2.png"},
	}

	return products, nil
}

func (dbh DevDatabaseHandler) GetCheckouts(ctx context.Context, userID int) ([]Checkout, error) {
	checkouts := []Checkout{
		{
			Product:         Product{Price: 100, Image: "image/product1.png"},
			ProductQuantity: 111,
		},
		{
			Product:         Product{Price: 200, Image: "image/product2.png"},
			ProductQuantity: 222,
		},
	}

	return checkouts, nil
}

func (dbh DevDatabaseHandler) CreateCheckout(userID int, productID int, productQuantity int) (string, error) {
	return "2022-03-04", nil
}
