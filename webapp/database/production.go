package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ProdDatabaseHandler struct {
	DB *sql.DB
}

func NewProdDatabaseHandler(db *sql.DB) ProdDatabaseHandler {
	return ProdDatabaseHandler{DB: db}
}

func (dbh ProdDatabaseHandler) InitDatabase() error {
	jsonFromFile, err := ioutil.ReadFile(InitDataJSONFileName)
	if err != nil {
		return err
	}

	var jsonData Blob
	if err := json.Unmarshal(jsonFromFile, &jsonData); err != nil {
		return err
	}

	db := dbh.DB

	queryTablenames := `
SELECT
  table_name
FROM
  information_schema.tables
WHERE table_name IN ('products', 'users', 'checkouts')
`

	tableNames, err := db.Query(queryTablenames)
	if err != nil {
		return err
	}

	existProducts, existUsers, existCheckouts := false, false, false

	for tableNames.Next() {
		table := ""
		tableNames.Scan(&table)

		switch strings.ToLower(table) {
		case "products":
			existProducts = true
		case "users":
			existUsers = true
		case "checkouts":
			existCheckouts = true
		}
	}

	// Don't use "IF EXISTS" as it is not supported by Spanner PGAdapter.
	queryDropProductsTables := "DROP TABLE products"
	if existProducts {
		if _, err := db.Exec(queryDropProductsTables); err != nil {
			return err
		}
	}

	queryDropUsersTables := "DROP TABLE users"
	if existUsers {
		if _, err := db.Exec(queryDropUsersTables); err != nil {
			return err
		}
	}

	queryDropCheckoutsTables := "DROP TABLE checkouts"
	if existCheckouts {
		if _, err := db.Exec(queryDropCheckoutsTables); err != nil {
			return err
		}
	}

	// Don't use "IF EXISTS" as it is not supported by Spanner PGAdapter.
	queryCreateProductsTable := `
	CREATE TABLE products (
		id INT64 NOT NULL,
		name STRING(20) NOT NULL,
		price INT64 NOT NULL,
		image STRING(100) NOT NULL
	) PRIMARY KEY(id)
	`

	queryCreateUsersTable := `
	CREATE TABLE users (
		id INT64 NOT NULL,
		name STRING(20) NOT NULL
	) PRIMARY KEY(id)
	`

	queryCreateCheckoutsTable := `
	CREATE TABLE checkouts (
		id STRING(40) NOT NULL,
		user_id INT64,
		product_id INT64,
		product_quantity INT64,
		created_at DATE
	) PRIMARY KEY(id)
	`

	if _, err := db.Exec(queryCreateProductsTable); err != nil {
		return err
	}

	if _, err := db.Exec(queryCreateUsersTable); err != nil {
		return err
	}

	if _, err := db.Exec(queryCreateCheckoutsTable); err != nil {
		return err
	}

	queryInsertProduct := "INSERT INTO products (id, name, price, image) VALUES(?, ?, ?, ?)"
	for _, product := range jsonData.Products {
		if _, err := db.Exec(queryInsertProduct, product.ID, product.Name, product.Price, product.Image); err != nil {
			return err
		}
	}

	queryInsertUser := "INSERT INTO users (id, name) VALUES(?, ?)"
	for _, user := range jsonData.Users {
		if _, err := db.Exec(queryInsertUser, user.ID, user.Name); err != nil {
			return err
		}
	}

	return nil
}

func (dbh ProdDatabaseHandler) GetProduct(id int) (Product, error) {
	var product Product

	db := dbh.DB
	query := "SELECT id, name, price, image FROM products WHERE id = ?"
	if err := db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Image); err != nil {
		return product, err
	}

	return product, nil
}

func (dbh ProdDatabaseHandler) GetProducts() ([]Product, error) {
	var products []Product

	db := dbh.DB
	query := "SELECT id, name, price, image FROM products"
	rows, err := db.Query(query)
	if err != nil {
		return products, err
	}

	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Image); err != nil {
			return products, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (dbh ProdDatabaseHandler) GetCheckouts(ctx context.Context, userID int) ([]Checkout, error) {
	var checkouts []Checkout

	db := dbh.DB
	query := `
	SELECT
	  products.name                           AS product_name,
	  products.image                          AS product_image,
	  checkouts.product_quantity              AS checkout_product_quantity,
	  FORMAT_DATE('%F', checkouts.created_at) AS checkout_created_at
	FROM checkouts
	INNER JOIN products ON checkouts.product_id = products.id
	WHERE checkouts.user_id = ?
	`

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return checkouts, err
	}

	for rows.Next() {
		var checkout Checkout
		if err := rows.Scan(&checkout.Product.Name, &checkout.Product.Image, &checkout.ProductQuantity, &checkout.CreatedAt); err != nil {
			return checkouts, err
		}

		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

func (dbh ProdDatabaseHandler) CreateCheckout(userID int, productID int, productQuantity int) (string, error) {
	uuidObj, err := uuid.NewRandom()
	createdAt := time.Now().Format("2006-01-02")
	checkoutID := uuidObj.String()
	if err != nil {
		return createdAt, nil
	}

	db := dbh.DB
	query := "INSERT INTO checkouts (id, user_id, product_id, product_quantity, created_at) VALUES (?, ?, ?, ?, ?)"
	if _, err := db.Exec(query, checkoutID, userID, productID, productQuantity, createdAt); err != nil {
		return createdAt, err
	}

	return createdAt, nil
}

func (dbh ProdDatabaseHandler) GetCheckout(checkoutID string) (Checkout, error) {
	checkout := Checkout{
		Product: Product{},
	}

	// err := dbh.Conn.Joins("Product").Find(&checkout, checkoutID).Error
	db := dbh.DB
	query := `
	SELECT
	  checkouts.id,
	  users.id,
	  users.name,
	  products.id,
	  products.name,
	  products.price,
	  products.image,
	  checkouts.product_quantity,
	  FORMAT_DATE('%F', checkouts.created_at)
	FROM checkouts
	LEFT JOIN users ON checkouts.user_id = users.id
	LEFT JOIN products ON checkouts.product_id = products.id
	WHERE checkouts.id = ?
	`
	if err := db.QueryRow(query, checkoutID).Scan(
		&checkout.Product.ID,
		&checkout.Product.Name,
		&checkout.Product.Price,
		&checkout.Product.Image,
		&checkout.ProductQuantity,
		&checkout.CreatedAt,
	); err != nil {
		return checkout, err
	}

	return checkout, nil
}
