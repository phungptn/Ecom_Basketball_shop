package database

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	sqlitecloud "github.com/sqlitecloud/sqlitecloud-go"
)

func QueryAllProduct(db *sqlitecloud.SQCloud, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	offsetStr := vars["offset"]

	// Convert offset to integer
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		http.Error(w, "Invalid offset value", http.StatusBadRequest)
		return err
	}

	sql := `SELECT s.productID, p.productName, s.size, s.stock, s.price
		FROM size s
		JOIN product p ON s.productID = p.productID
		WHERE s.stock > 0
		LIMIT 10 OFFSET ?`

	values := []interface{}{offset * 10}
	rows, err := db.SelectArray(sql, values)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(rows.ToJSON()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

func QueryAllOrders(db *sqlitecloud.SQCloud, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	method := vars["method"]
	offsetStr := vars["offset"]

	// Convert offset to integer
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		http.Error(w, "Invalid offset value", http.StatusBadRequest)
		return err
	}

	var sql string
	var values []interface{}

	// Determine the SQL query based on the method
	switch method {
	case "prepared", "sent", "processed", "finished":
		sql = `SELECT o.orderID, o.userID, o.date, o.shippingAddress, o.billingAddress, o.price, o.status, o.payStatus, o.method,
		              ud.fullName, COALESCE(ud.phoneNumber, 'null') AS phoneNumber
		       FROM orders o
		       JOIN userDetail ud ON o.userID = ud.userID
		       WHERE o.status = ?
			   ORDER BY o.date DESC
		       LIMIT 5 OFFSET ?`
		values = []interface{}{method, offset * 5}
	case "paid", "unpaid":
		sql = `SELECT o.orderID, o.userID, o.date, o.shippingAddress, o.billingAddress, o.price, o.status, o.payStatus, o.method,
		              ud.fullName, COALESCE(ud.phoneNumber, 'null') AS phoneNumber
		       FROM orders o
		       JOIN userDetail ud ON o.userID = ud.userID
		       WHERE o.payStatus = ?
			   ORDER BY o.date DESC
		       LIMIT 5 OFFSET ?`
		values = []interface{}{method, offset * 5}
	case "all":
		sql = `SELECT o.orderID, o.userID, o.date, o.shippingAddress, o.billingAddress, o.price, o.status, o.payStatus, o.method,
		              ud.fullName, COALESCE(ud.phoneNumber, 'null') AS phoneNumber
		       FROM orders o
		       JOIN userDetail ud ON o.userID = ud.userID
			   ORDER BY o.date DESC
		       LIMIT 5 OFFSET ?`
		values = []interface{}{offset * 5}
	default:
		http.Error(w, "Invalid method value", http.StatusBadRequest)
		return fmt.Errorf("Invalid method value")
	}

	// Execute the query
	rows, err := db.SelectArray(sql, values)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(rows.ToJSON()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func QueryItemsByOrderID(db *sqlitecloud.SQCloud, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	orderID := vars["orderID"]
	sql := `
        SELECT oi.orderID, oi.productID, p.productName, oi.size, oi.quantity, oi.price
        FROM orderItem oi
        JOIN product p ON oi.productID = p.productID
        WHERE oi.orderID = ?
    `

	values := []interface{}{orderID}
	rows, err := db.SelectArray(sql, values)
	if err != nil {
		return err
	}

	// Convert the result rows to JSON and write to the response
	_, err = w.Write([]byte(rows.ToJSON()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func QueryAllCustomers(db *sqlitecloud.SQCloud, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	offsetStr := vars["offset"]

	// Convert offset to integer
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		http.Error(w, "Invalid offset value", http.StatusBadRequest)
		return err
	}

	// SQL query to retrieve customer information where the role is 'customer'
	sql := `SELECT u.email, ud.fullName, COALESCE(ud.phoneNumber, 'null') AS phoneNumber, 
                   COALESCE(ud.dob, 'null') AS dob, COALESCE(ud.memberSince, 'null') AS memberSince 
            FROM users u
            JOIN userDetail ud ON u.userID = ud.userID
            WHERE u.role = 'customer'
            LIMIT 5 OFFSET ?`

	values := []interface{}{offset * 5}

	// Execute the query
	rows, err := db.SelectArray(sql, values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	_, err = w.Write([]byte(rows.ToJSON()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
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

func DeleteProduct(db *sqlitecloud.SQCloud, productID string) error {
	// SQL statement to update all sizes of the specified product to set stock to 0
	updateStockSQL := "UPDATE size SET stock = 0 WHERE productID = ?"

	// Execute the SQL update statement
	err := db.ExecuteArray(updateStockSQL, []interface{}{productID})
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Optionally, log or handle the success case
	fmt.Println("All sizes of the product have been set to 0 stock successfully.")
	return nil
}

func UpdateOrderStatus(db *sqlitecloud.SQCloud, orderID string, status string) error {
	// SQL statement to update all sizes of the specified product to set stock to 0
	updateStockSQL := "UPDATE orders SET status = ? WHERE orderID = ?"

	// Execute the SQL update statement
	err := db.ExecuteArray(updateStockSQL, []interface{}{status, orderID})
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Optionally, log or handle the success case
	fmt.Println("Update order status successfully")
	return nil
}

type StatResponse struct {
	TotalOrder      int     `json:"totalOrder"`
	TotalProduct    int     `json:"totalProduct"`
	UnfinishedOrder int     `json:"unfinishedOrder"`
	FinishedOrder   int     `json:"finishedOrder"`
	TotalMonthMoney float64 `json:"totalMonthMoney"`
	TotalMoney      float64 `json:"totalMoney"`
	TotalCustomer   int     `json:"totalCustomer"`
	TopProduct      struct {
		ProductName   string  `json:"productName"`
		TotalQuantity int     `json:"totalQuantity"`
		TotalPrice    float64 `json:"totalPrice"`
	} `json:"topProduct"`
}

func QueryStat(db *sqlitecloud.SQCloud, w http.ResponseWriter) error {
	var result StatResponse

	// Total orders
	rows, err := db.Select(`SELECT count(orderID) as totalOrder FROM "orders";`)
	if err != nil {
		return err
	}
	result.TotalOrder = int(rows.GetInt32Value_(0, 0))

	// Total products
	rows, err = db.Select(`SELECT count(productID) AS totalProduct FROM "product";`)
	if err != nil {
		return err
	}
	result.TotalProduct = int(rows.GetInt32Value_(0, 0))

	// Unfinished orders
	rows, err = db.Select(`SELECT count(orderID) AS unfinishedOrder FROM "orders" WHERE status != 'finished';`)
	if err != nil {
		return err
	}
	result.UnfinishedOrder = int(rows.GetInt32Value_(0, 0))

	// Finished orders
	rows, err = db.Select(`SELECT count(orderID) AS finishedOrder FROM "orders" WHERE status = 'finished';`)
	if err != nil {
		return err
	}
	result.FinishedOrder = int(rows.GetInt32Value_(0, 0))

	// Total money this month
	rows, err = db.Select(`SELECT SUM(price) AS totalMonthMoney FROM orders WHERE strftime('%Y-%m', date) = strftime('%Y-%m', 'now');`)
	if err != nil {
		return err
	}
	result.TotalMonthMoney = rows.GetFloat64Value_(0, 0)

	// Total money
	rows, err = db.Select(`SELECT SUM(price) AS totalMoney FROM "orders";`)
	if err != nil {
		return err
	}
	result.TotalMoney = rows.GetFloat64Value_(0, 0)

	// Total customers
	rows, err = db.Select(`SELECT COUNT(userID) AS totalCustomer FROM "users" WHERE role = 'customer';`)
	if err != nil {
		return err
	}
	result.TotalCustomer = int(rows.GetInt32Value_(0, 0))

	// Top product
	rows, err = db.Select(`
		SELECT p.productName,
			SUM(oi.quantity) AS totalQuantity,
			SUM(oi.quantity * oi.price) AS totalPrice
		FROM orderItem oi
		JOIN product p ON oi.productID = p.productID
		GROUP BY p.productName
		ORDER BY totalQuantity DESC
		LIMIT 1;`)
	if err != nil {
		return err
	}

	result.TopProduct.ProductName = rows.GetStringValue_(0, 0)
	result.TopProduct.TotalQuantity = int(rows.GetInt32Value_(0, 1))
	result.TopProduct.TotalPrice = rows.GetFloat64Value_(0, 2)

	// Set headers and encode the response
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}
