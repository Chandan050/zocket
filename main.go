package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()
	route.HandleFunc("/product", CreateProducts).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", route))
}

func CreateProducts(w http.ResponseWriter, r *http.Request) {
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Fatalf("[error] while decoding the product %f", err)
	}

}
