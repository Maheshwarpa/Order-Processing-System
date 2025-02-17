package Consumer

import (
	"OPS/module/DatabaseConn"
	"OPS/module/Event"
	"OPS/module/Notification"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
	//"github.com/jackc/pgx/v5/pgxpool"
)

var ps Event.ProcessingResponse

// ConsumeMessages starts the Kafka consumer
func ConsumeMessages(brokers []string, topic string, groupID string, wg *sync.WaitGroup) {

	//fmt.Println("Entered ConsumerMessages")
	// Configure Kafka consumer
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create a consumer handler
	handler := ConsumerHandler{}
	handler.wg = wg

	// Create consumer group
	for {
		consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
		if err != nil {
			log.Fatalf("Error creating consumer group: %v", err)
		}
		defer consumerGroup.Close()

		err1 := consumerGroup.Consume(context.Background(), []string{topic}, &handler)
		if err1 != nil {
			log.Fatalf("Error consuming messages: %v", err1)
		}

	}

}

// ConsumerHandler handles consumed messages
type ConsumerHandler struct {
	wg *sync.WaitGroup
}

// Setup is run at the beginning of a new session
func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup is run at the end of a session
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from Kafka
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	//fmt.Println("Entered  ConsumeClaim Section")
	for msg := range claim.Messages() {
		//fmt.Println("Entered claim message")
		var event Event.OrderCreatedResponse
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}
		// dbPool, _ := DatabaseConn.ConnectDB()
		// if err != nil {
		// 	log.Fatalf("Failed to connect to the database: %v", err)
		// }
		//defer dbPool.Close()
		price, err := DatabaseConn.GetPrice(DatabaseConn.DbPool, event.Order_id)
		if err != nil {
			return fmt.Errorf("unable to fetch the price from the database order based on userid %v", err)

		}
		if price >= 1000 {
			ps.Order_id = event.Order_id
			ps.Payment_status = "FAILED"

		} else {
			ps.Order_id = event.Order_id
			ps.Payment_status = "SUCCESS"
		}

		if Event.PRList == nil {
			Event.PRList = []Event.ProcessingResponse{}
		}
		Event.PRList = append(Event.PRList, ps)
		err = DatabaseConn.AddPaymentStatus(DatabaseConn.DbPool, ps)
		if err != nil {
			log.Fatalf("error in adding data in paymentstatus table %v", err)
		}
		jsondata, err := json.Marshal(ps)
		if err != nil {
			fmt.Errorf("Unable to marshal the payment status data")
		}
		fmt.Printf("PaymentProcessedEvent status:\n%+v\n", string(jsondata))
		//fmt.Printf("Received message: OrderID=%s, Status=%s\n", event.Order_id, ps.Payment_status)
		notificationResponse, err := Notification.GetNotification(ps)
		if err != nil {
			log.Printf("Error sending notification: %v", err)
		} else {
			fmt.Printf("Notification Response: \n%+v\n", notificationResponse)
		}

		h.wg.Done()

		session.MarkMessage(msg, "")

	}
	//fmt.Println("Outside of claim message")
	return nil
}
