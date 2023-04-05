package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
	"time"
)

const (
	groupID = "bitcoin-1"
	topic   = "bitcoin"
)

var (
	bootstrapServers = os.Getenv("BOOTSTRAP_SERVERS")
)

func main() {
	c := NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	defer c.Close()

	err := c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topics: %v\n", err)
		return
	}

	for {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func NewConsumer(config *kafka.ConfigMap) *kafka.Consumer {
	c, err := kafka.NewConsumer(config)
	if err != nil {
		panic(err)
	}
	return c
}
