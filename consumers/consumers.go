package consumers

import (
	"ZocketAssigment"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nfnt/resize"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func startKafkaConsumer() {
	//creating kafka consume and adding kafka consume config
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092", //Kafka broker address
		"group.id":          "my-consumer-group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Subscribe to the Kafka topic
	topic := "test-topic" //Kafka topic
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to Kafka topic: %v", err)
	}

	// Consume Kafka messages
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			// Process the message (e.g., download images and update the database)
			processKafkaMessage(msg.Value)
		} else {
			log.Printf("Error reading Kafka message: %v", err)
		}
	}
}

var product *ZocketAssigment.Product

func GetProductData(data ZocketAssigment.Product) {
	product := &data

}

func processKafkaMessage(message []byte) {
	// Implement message processing logic here
	productID := string(message)
	log.Printf("Received message from Kafka: Product ID %s", productID)
	resizedImages := downloadimages(product)

	dsn := "user:password@tcp(localhost:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// Handle the error
	}

	if err := db.First(&product, product.productID).Error; err != nil {
		// Handle the error (e.g., product not found)
	}

	// Update the field you want to change
	product.Name = newProductName

	// Save the changes back to the database
	if err := db.Save(&product).Error; err != nil {
		// Handle the error
	}

}

func downloadimages(i []string) []string {

	folderPath := "./resizedPhoto"
	Width := uint(300)  //  width
	Height := uint(200) //  height

	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating file:", err)
	}

	var imageAddresses []string

	for _, url := range i {

		// Send an HTTP GET request to the URL to download the image
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error creating file: %f", err)
		}

		defer resp.Body.Close()

		// Extract the image filename from the URL
		imageFilename := url[strings.LastIndex(url, "/")+1:]

		// Create a full file path including the folder path and image filename
		fullPath := folderPath + "/" + imageFilename

		// Create a new file to save the downloaded image
		file, err := os.Create(fullPath)
		if err != nil {
			log.Fatalf("Error creating file: %f", err)
		}

		defer file.Close()

		// Decode the downloaded image
		img, _, err := image.Decode(resp.Body)
		if err != nil {
			log.Fatalf("Error creating file: %f", err)
		}

		// Resize the image to the specified width and height
		resizedImg := resize.Resize(Width, Height, img, resize.Lanczos3)

		// Encode the resized image and save it to the file
		err = jpeg.Encode(file, resizedImg, nil)
		if err != nil {
			log.Fatalf("Error creating file: %f", err)
		}

		fmt.Printf("Downloaded and resized image from %s and saved as %s\n", url, fullPath)
		imageAddresses = append(imageAddresses, fullPath)
	}
	return imageAddresses
}
