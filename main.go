package main

import (
	"context"
	"encoding/json"
	"fmt"
	schema "kafka/database"
	producer "kafka/producers"
	"log"
	"net/http"

	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func CreateProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var product schema.Product

	// var user schema.User

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}
	fmt.Println("encodeed")
	// Connect to MongoDB
	_, collection, _ := schema.ConnectToMongo()

	fmt.Println("conected to users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// Insert the product into the MongoDB collection
	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to insert product into the database", http.StatusInternalServerError)
		return
	}
	// filter := bson.M{"ID": product.UserID}

	// // Perform the findOne query
	// var userfilter bson.M
	// err = collection.FindOne(ctx, filter).Decode(&userfilter)
	// if err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		fmt.Println("User not found")
	// 	} else {
	// 		log.Fatal(err)
	// 	}
	// } else {
	// 	Usersresult, err := collectionUsers.InsertOne(ctx, user)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		http.Error(w, "Failed to insert Users into the database ", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	fmt.Println("user not found so created ", Usersresult)

	// }

	response := map[string]interface{}{
		"message":   "Product created successfully",
		"productID": result.InsertedID,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	producer.SendMessageToKafka(result.InsertedID)

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func main() {

	// Create a new HTTP router
	r := mux.NewRouter()

	// Define your API routes
	r.HandleFunc("/products", CreateProduct).Methods("POST")

	// Start the HTTP server
	port := 8080
	serverAddr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server listening on %s\n", serverAddr)
	http.Handle("/", r)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal(err)
	}

}
