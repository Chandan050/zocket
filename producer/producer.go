package producer

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func SendMessageToKafka(producer *kafka.Producer, productID int) {
	// Produce a message to the Kafka topic
	topic := "test-topic" //  Kafka topic
	message := fmt.Sprintf("%d", productID)

	// Produce a message to Kafka
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
	} else {
		fmt.Println(message)
	}
	producer.Flush(15 * 100)
	producer.Close()

}
