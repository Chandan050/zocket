package api

import (
	"context"
	"encoding/json"
	d "kafka/database"
	p "kafka/producers"
	"net/http"
	"time"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var product d.Product
	// var user d.User

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	inserted, err := d.ProductsCollection.InsertOne(context.Background(), product)

	if err != nil {
		http.Error(w, "Failed to insert ", http.StatusInternalServerError)

		return
	}
	p.SendMessageToKafka(inserted.InsertedID)

	// d.USersCollection.InsertOne(context.Background(),user)

	response := map[string]interface{}{"message": "Product created successfully", "product": inserted.InsertedID}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to serialize JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
