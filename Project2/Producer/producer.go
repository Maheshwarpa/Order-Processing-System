package Producer

import (
	"OPS/module/Event"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

//const topic = "orders"

func PublishMessage(brokers []string, topic string, message *Event.OrderCreatedResponse, wg *sync.WaitGroup) {
	// Configure Kafka producer
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Create a new Kafka producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}
	defer producer.Close()

	// Convert struct to JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to serialize message: %v", err)
	}

	// Create Kafka message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(jsonMessage),
	}

	// Send message to Kafka
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	wg.Done()
}
