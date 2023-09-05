package main

import (
	"context"
	"fmt"
	p "kafka/producers"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockMongoCollection is a mock implementation of the MongoDB collection.
type MockMongoCollection struct {
	data map[string]interface{}
}

func (m *MockMongoCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	id := filter.(bson.M)["_id"].(int)
	result := m.data[fmt.Sprintf("%d", id)]
	return &mongo.SingleResult{
		Err:  nil,
		SRVD: result,
	}
}

func (m *MockMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return &mongo.UpdateResult{}, nil
}

func TestDownloadCompressAndStoreImage(t *testing.T) {
	imageURL := "https://example.com/image.jpg"
	outputFilePath, err := downloadCompressAndStoreImage(imageURL)
	assert.NoError(t, err)
	assert.NotEmpty(t, outputFilePath)
	// Add assertions to check if the image was downloaded, compressed, and stored correctly.
}

func TestProcessProductImage(t *testing.T) {
	// Create a mock MongoDB collection
	mockCollection := &MockMongoCollection{
		data: map[string]interface{}{
			"42": map[string]interface{}{
				"_id": 42,
				"product_images": []string{
					"https://example.com/image1.jpg",
					"https://example.com/image2.jpg",
				},
			},
		},
	}

	// Call the function being tested
	processProductimage(42, mockCollection)

	// Add assertions to check if the product was updated with the correct compressed image paths.
}

func TestMainFunction(t *testing.T) {
	// Replace the actual implementation of ReceiveLatestFomKafka with a mock for testing purposes.
	p.ReceiveLatestFomKafka = func() (int, error) {
		return 42, nil
	}

	// Create a mock MongoDB collection
	mockCollection := &MockMongoCollection{
		data: map[string]interface{}{
			"42": map[string]interface{}{
				"_id": 42,
				"product_images": []string{
					"https://example.com/image1.jpg",
					"https://example.com/image2.jpg",
				},
			},
		},
	}

	// Call the main function
	main()

	// Add assertions to check if the main function behaves as expected.
}

func TestGetImageFileName(t *testing.T) {
	imageURL := "https://example.com/images/image.jpg"
	fileName := getImageFileName(imageURL)
	assert.Equal(t, "image.jpg", fileName)
}
