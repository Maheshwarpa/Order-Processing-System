package Notification

import (
	"OPS/module/Event"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
)

type NotificationResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

func GetNotification(ps Event.ProcessingResponse) (string, error) {
	//fmt.Println("Notification is being called")
	input := ps.Order_id
	parts := strings.Split(input, "_") // Split by "_"
	lastPart := parts[len(parts)-1]

	var (
		ctx       = context.Background()
		user_id   = lastPart
		redisAddr = "localhost:6379"
		message   = ""
	)

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Store a notification
	if ps.Payment_status == "SUCCESS" {
		message = "Your payment was successful for " + ps.Order_id
		//message += ps.Order_id
	} else {
		message = "Your payment was failed for " + ps.Order_id
		//message += ps.Order_id
	}
	//notification := "New message from Alice!"
	err := client.LPush(ctx, user_id, message).Err()
	if err != nil {
		log.Fatalf("Error storing notification: %v", err)
	}

	// Keep only the latest 100 notifications
	client.LTrim(ctx, user_id, 0, 99)

	// Create response
	response := NotificationResponse{
		UserID:  user_id,
		Message: message,
	}

	// Convert response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error creating JSON response: %v", err)
		return "", err
	}

	fmt.Println("Notification stored successfully!")
	return string(jsonResponse), nil

}
