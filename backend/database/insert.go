package database

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	sqlitecloud "github.com/sqlitecloud/sqlitecloud-go"
)

func recordExists(db *sqlitecloud.SQCloud, tableName string, attribute string, id string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE %s = '%s'", tableName, attribute, id)
	rows, err := db.Select(query)
	if err != nil && rows.GetInt32Value_(0, 0) == 0 {
		return false, err
	}
	return rows.GetInt32Value_(0, 0) != 0, nil
}

func InsertUser(db *sqlitecloud.SQCloud, email string, password string) error {
	var id string
	for {
		id = uuid.New().String()
		exists, err := recordExists(db, "users", "userID", id)
		if err != nil {
			return err
		}
		if !exists {
			break
		}
	}
	insertUserSQL := "INSERT INTO users (userID, email, password) VALUES (?, ?, ?)"
	values := []interface{}{id, email, password}

	err := db.ExecuteArray(insertUserSQL, values)
	return err
}

func InsertCategory(db *sqlitecloud.SQCloud, name string, description string) error {
	var id string
	for {
		id = uuid.New().String()
		exists, err := recordExists(db, "category", "categoryID", id)
		if err != nil {
			return err
		}
		if !exists {
			break
		}
	}

	insertCategorySQL := "INSERT INTO category (categoryID, categoryName, description) VALUES (?, ?, ?)"
	values := []interface{}{id, name, description}

	err := db.ExecuteArray(insertCategorySQL, values)
	return err
}

func InsertProduct(db *sqlitecloud.SQCloud, categoryID string, name string, description string, brand string, price string, stock string, dateAdded string, size string) error {
	var id string
	for {
		id = uuid.New().String()
		exists, err := recordExists(db, "product", "productID", id)
		if err != nil {
			return err
		}
		if !exists {
			break
		}
	}

	priceValue, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	stockInt, err := strconv.Atoi(stock)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	insertProductSQL := "INSERT INTO product (productID, categoryID, productName, description, brand, dateAdded) VALUES (?, ?, ?, ?, ?, ?)"
	values := []interface{}{id, categoryID, name, description, brand, dateAdded}
	err = db.ExecuteArray(insertProductSQL, values)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	insertSizeSQL := "INSERT INTO size (productID, size, stock, price) VALUES (?, ?, ?, ?)"
	values2 := []interface{}{id, size, stockInt, priceValue}
	err = db.ExecuteArray(insertSizeSQL, values2)
	return err
}
