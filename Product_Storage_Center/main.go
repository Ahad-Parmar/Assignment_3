package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Item-Product

type Product struct {
	UID   string  `json:"UID"`
	Name  string  `json:"Name"`
	Desc  string  `json:"Desc"`
	Price float64 `json:"Price"`
}

var items []Product

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Fprintf(w, "Endpoints called : HomePage()")
}

func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Function called: getitems()")

	json.NewEncoder(w).Encode(items)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/join")
	fmt.Println("Function called: createProduct()")

	var product Product
	_ = json.NewDecoder(r.Body).Decode(&product)

	items = append(items, product)

	json.NewEncoder(w).Encode(product)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/join")
	fmt.Println("Function called: deleteProduct()")

	params := mux.Vars(r)

	_deleteProductAtUid(params["uid"])

	json.NewEncoder(w).Encode(items)
}

func _deleteProductAtUid(uid string) {
	for index, product := range items {
		if product.UID == uid {
			
	  			// delete product from slice
		
			items = append(items[:index], items[index+1:]...)
			break
		}
	}
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/join")
	fmt.Println("Function called: updateProduct()")

	var product Product
	_ = json.NewDecoder(r.Body).Decode(&product)

	params := mux.Vars(r)

	// delete the item at UID

	_deleteProductAtUid(params["uid"])

	// create it with new data

	items = append(items, product)

	json.NewEncoder(w).Encode(items)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/items", getItems).Methods("GET")
	router.HandleFunc("/items", createProduct).Methods("POST")
	router.HandleFunc("/items/{uid}", deleteProduct).Methods("DELETE")
	router.HandleFunc("/items/{uid}", updateProduct).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	items = append(items, Product{
		UID:   "1",
		Name:  "Cellphone",
		Desc:  "Samsung S9 plus",
		Price: 65000,
	})
	items = append(items, Product{
		UID:   "2",
		Name:  "Cellphone",
		Desc:  "Samsung S10 plus",
		Price: 75000,
	})
	handleRequests()
}
