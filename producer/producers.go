package producers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func SendMessageToKafka(productID interface{}) {

	// Produce a message to the Kafka topic
	topic := "MicroserviceTopic" //  Kafka topic

	message := fmt.Sprintf("%d", productID)

	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"acks":              "all",
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatal("Failed to create admin client:", err)
	}
	defer producer.Close()

	// Produce a message to Kafka
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
	}

	producer.Flush(15 * 1000)

	fmt.Println("producerc sent msg to consumer")

}

func ReceiveLatestFomKafka() (int, error) {
	topic := "MicroserviceTopic"

	// Kafka consumer configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "my-consumer-group",
		"auto.offset.reset": "latest",
	}
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return 0, err
	}
	defer consumer.Close()

	// Subscribe to the topic
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return 0, err
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		fmt.Println(msg)
		if err == nil {
			int1, _ := strconv.Atoi(string(msg.Value))
			return int1, err
		} else {
			return 0, err
		}
	}
}
