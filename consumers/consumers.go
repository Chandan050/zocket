package main

import (
	"context"
	"fmt"
	"io"
	d "kafka/database"
	p "kafka/producers"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func downloadCompressAndStoreImage(url string, res int, i int) (string, error) {
	// Send an HTTP GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP GET request failed with status: %s", resp.Status)
	}

	// Create a directory to store the images if it doesn't exist
	err = os.MkdirAll("images", os.ModePerm)
	if err != nil {
		return "", err
	}

	// Create a prodcts folder  to store the compressed image
	imageFolder := fmt.Sprintf("product_%d", res)
	imageFoladerlocation := fmt.Sprintf("images/%s", imageFolder)

	err = os.MkdirAll(imageFoladerlocation, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Create a prodcts folder  to store the compressed image
	imagefileName := fmt.Sprintf("%s/images_%v_%v", imageFoladerlocation, res, i)
	outfilePath, err := os.Create(imagefileName + ".jpeg")
	if err != nil {
		return "", err
	}
	defer outfilePath.Close()
	// Copy the response body (image data) to the file
	_, err = io.Copy(outfilePath, resp.Body)
	if err != nil {
		return "", err
	}

	// Update the product with the compressed image paths
	return imagefileName, nil
}

var product d.Product

func main() {

	fmt.Println("in consumer")
	res, err := p.ReceiveLatestFomKafka()
	fmt.Println(res)
	if err != nil {
		log.Fatal("[error] while getting msg from producer", err)
	}

	filter := bson.M{"_id": res}
	_, ProductsCollection, _ := d.ConnectToMongo()
	if err := ProductsCollection.FindOne(context.Background(), filter).Decode(&product); err != nil {
		log.Fatal(err)
	}
	processProductimage(res, ProductsCollection)

}

func processProductimage(a int, ProductsCollection *mongo.Collection) {

	var images []string
	for _, v := range product.ProductImages {
		images = append(images, v)
	}

	var filepath []string

	for i, url := range images {
		outputFilePath, err := downloadCompressAndStoreImage(url, a, i)
		if err != nil {
			fmt.Printf("Error processing image: %v\n", err)
			continue
		}

		filepath = append(filepath, `D:\go\src\kafka1\consumers`+outputFilePath)

	}

	update := bson.M{"$set": bson.M{"compressed_product_images": filepath}}
	_, err := ProductsCollection.UpdateOne(context.Background(), bson.M{"_id": a}, update)
	if err != nil {
		fmt.Println(err)
	}

}
