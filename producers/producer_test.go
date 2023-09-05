package producers

import (
	"fmt"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/mock"
)

// MockKafkaProducer is a mock implementation of the Kafka producer.
type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockKafkaProducer) Flush(timeoutMs int) int {
	args := m.Called(timeoutMs)
	return args.Int(0)
}

func (m *MockKafkaProducer) Close() {
	m.Called()
}

func TestSendMessageToKafka(t *testing.T) {
	// Create a mock Kafka producer
	mockProducer := new(MockKafkaProducer)

	// Replace the real Kafka producer with the mock
	originalProducer := kafkaProducer
	kafkaProducer = mockProducer
	defer func() {
		kafkaProducer = originalProducer
	}()

	// Define the test input
	productID := 42

	// Set up expectations for the mock producer
	expectedMessage := fmt.Sprintf("%d", productID)
	mockProducer.On("Produce", mock.Anything, mock.Anything).Return(nil)
	mockProducer.On("Flush", 15000).Return(0)
	mockProducer.On("Close").Return()

	// Call the function being tested
	SendMessageToKafka(productID)

	// Assert that the Produce method was called with the expected message
	mockProducer.AssertCalled(t, "Produce", &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: new(string), Partition: kafka.PartitionAny},
		Value:          []byte(expectedMessage),
	}, mock.Anything)

	// Assert that the Flush and Close methods were called
	mockProducer.AssertCalled(t, "Flush", 15000)
	mockProducer.AssertCalled(t, "Close")
}
