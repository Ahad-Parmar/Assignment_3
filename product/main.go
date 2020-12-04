package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Order represents the model for an order

type Order struct {             // Default table name will be `orders`
	// gorm.Model
	OrderID      uint   `json:"orderid" gorm:"primary_key"`
	CustomerName string `json:"customername"`
	Items        []Item `json:"items" gorm:"foreignkey:OrderID"`
}



// Item represents the model for an item in the order
type Item struct {
	// gorm.Model
	LineItemID  uint   `json:"lineitemid" gorm:"primary_key"`
	ItemCode    string `json:"itemcode"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	OrderID     uint   `json:"-"`
}

var db *gorm.DB

func initDB() {
	var err error
	DataSourceName := "root:password@tcp(localhost:3306)/?parseTime=True"
	db, err = gorm.Open("mysql", DataSourceName)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect with database")
	}

	// Create the database. This is a one-time step.


	// db.Exec("CREATE DATABASE order_db")

	db.Exec("USE order_db")

	// Migration to create tables for Order and Item schema

	db.AutoMigrate(&Order{}, &Item{})
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	json.NewDecoder(r.Body).Decode(&order)
	// Creates new order by inserting records in the `orders` and `items` table
	db.Create(&order)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var orders []Order
	db.Preload("Items").Find(&orders)
	json.NewEncoder(w).Encode(orders)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	InputOrderID := params["orderid"]

	var order Order
	db.Preload("Items").First(&order, InputOrderID)
	json.NewEncoder(w).Encode(order)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var UpdatedOrder Order
	json.NewDecoder(r.Body).Decode(&UpdatedOrder)
	db.Save(&UpdatedOrder)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UpdatedOrder)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	InputOrderID := params["orderid"]

	// Convert `orderId` string param to uint64 - important step

	id64, _ := strconv.ParseUint(InputOrderID, 10, 64)

	// Convert uint64 to uint

	IdToDelete := uint(id64)

	db.Where("order_id = ?", IdToDelete).Delete(&Item{})
	db.Where("order_id = ?", IdToDelete).Delete(&Order{})
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	router := mux.NewRouter()
	// Create
	router.HandleFunc("/orders", CreateOrder).Methods("POST")
	// Read
	router.HandleFunc("/orders/{orderid}", GetOrder).Methods("GET")
	// Read-everything
	router.HandleFunc("/orders", GetOrders).Methods("GET")
	// Update
	router.HandleFunc("/orders/{orderid}", UpdateOrder).Methods("PUT")
	// Delete
	router.HandleFunc("/orders/{orderid}", DeleteOrder).Methods("DELETE")
	// Initialize db connection
	initDB()

	log.Fatal(http.ListenAndServe(":8080", router))
}
