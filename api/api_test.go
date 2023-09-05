package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kafka/api"
	d "kafka/database"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
// MockMongoCollection is a mock implementation of the MongoDB collection.
type MockMongoCollection struct {
	data map[string]interface{}
}

func (m *MockMongoCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	// Simulate the insertion of a document and return a mock result
	return &mongo.InsertOneResult{InsertedID: "mock-inserted-id"}, nil
}

func TestCreateProductWithMongoDB(t *testing.T) {
	// Create a test product to send in the request
	product := d.Product{
		UserID:               1,
		ProductName:          "Sample Product",
		ProductDescription:   "This is a sample product.",
		ProductImages:        []string{"image1.jpg", "image2.jpg"},
		ProductPrice:         19.99,
		CompressedImagesPath: nil,
		CreatedAt:            time.Time{},
		UpdatedAt:            time.Time{},
	}
	// / Serialize the product to JSON
	productJSON, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("Failed to serialize product to JSON: %v", err)
	}

	// Create a request with the JSON payload
	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(productJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a test HTTP response recorder
	recorder := httptest.NewRecorder()

	// Call the CreateProduct handler function
	CreateProduct(recorder, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Unmarshal the response JSON
	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
	// Check the response message and inserted ID
	assert.Equal(t, "Product created successfully", response["message"])
	assert.Equal(t, "mock-inserted-id", response["product"])

	// Verify that the product was inserted into the mock MongoDB collection
	assert.Equal(t, 1, len(mockCollection.data))

	// Verify that the product details match the inserted product
	insertedProduct := mockCollection.data["mock-inserted-id"].(map[string]interface{})
	assert.Equal(t, product.ID, int(insertedProduct["_id"].(int32)))
	assert.Equal(t, product.UserID, int(insertedProduct["user_id"].(int32)))
	assert.Equal(t, product.ProductName, insertedProduct["product_name"].(string))
	assert.Equal(t, product.ProductDescription, insertedProduct["product_description"].(string))
	assert.Equal(t, product.ProductImages, insertedProduct["product_images"].([]string))
	assert.Equal(t, product.ProductPrice, insertedProduct["product_price"].(float64))
	// You can add more assertions based on your requirements

	
}

	
\


	

	// Create a response recorder to capture the response
	// Set up MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Error creating MongoDB client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	

	// Parse the response JSON
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify the response message
	expectedMessage := "Product created successfully"
	assert.Equal(t, expectedMessage, response["message"])

	// Verify that the product ID is present in the response
	assert.NotNil(t, response["product"])

	// Clean up: Remove the test product from the database
	productID, ok := response["product"].(string)
	if !ok {
		t.Fatalf("Failed to extract product ID from the response")
	}

	// Convert the product ID to an ObjectId
	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		t.Fatalf("Failed to convert product ID to ObjectId: %v", err)
	}

	// Delete the test product from the database
	_, err = d.ProductsCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		t.Fatalf("Failed to delete the test product: %v", err)
	}
