package main

import (
	"ZocketAssignment/producer"
	"ZocketAssignment/consumers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	// Create a new producer and start it in background
	route.HandleFunc("/product", CreateProducts).Methods(http.MethodPost)
	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", route))
}

func CreateProducts(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	var product Product
	consumers.SendProductData(product)
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Fatalf("[error] while decoding the product %f", err)
	}

	response := map[string]interface{}{
		"message":   "Product created successfully",
		"productID": product.ProductID,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

	msgProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092", //Kafka broker address
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer msgProducer.Close()

	producer.SendMessageToKafka(msgProducer, int(product.ProductID))

}
