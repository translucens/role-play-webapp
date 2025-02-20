package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
)

type DatabaseHandler interface {
	InitDatabase() error
	GetProduct(id int) (Product, error)
	GetProducts() ([]Product, error)
	GetCheckouts(ctx context.Context, userID int) ([]Checkout, error)
	CreateCheckout(userID int, productID int, productQuantity int) (string, error)
}

const InitDataJSONFileName = "initdata.json"

type Database struct {
	DB *sql.DB
}

type Product struct {
	ID    int
	Name  string
	Price int
	Image string
}

type User struct {
	ID   int
	Name string
}

type Checkout struct {
	Product         Product
	ProductQuantity int
	CreatedAt       string
}

type Blob struct {
	Products []Product `json:"products"`
	Users    []User    `json:"users"`
}

func NewDatabaseHandler(environment string, db *sql.DB) (DatabaseHandler, error) {
	switch environment {
	case "production":
		return NewProdDatabaseHandler(db), nil
	case "development":
		return NewDevDatabaseHandler(db), nil
	default:
		return nil, fmt.Errorf("environment: %s is not supported", environment)
	}
}

func InitializeDevDBConn() (*sql.DB, error) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	return mockDB, nil
}
