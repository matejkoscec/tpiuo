package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewProducer(config *kafka.ConfigMap) *kafka.Producer {
	p, err := kafka.NewProducer(config)
	if err != nil {
		panic(err)
	}
	return p
}

func CreateTopics(a *kafka.AdminClient, topicSpecification []kafka.TopicSpecification) {
	if a == nil {
		fmt.Println("Cannot create topics, AdminClient is nil")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		panic("ParseDuration(60s)")
	}
	results, err := a.CreateTopics(ctx, topicSpecification, kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		fmt.Printf("Failed to create topic: %v\n", err)
		os.Exit(1)
	}

	for _, result := range results {
		fmt.Printf("%s\n", result)
	}
}

func NewAdminClient(config *kafka.ConfigMap) *kafka.AdminClient {
	a, err := kafka.NewAdminClient(config)
	if err != nil {
		fmt.Printf("Failed to create Admin client: %s\n", err)
		return nil
	}
	return a
}

func Produce(producer *kafka.Producer, messages []Message, topic string) {
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	for _, m := range messages {
		for _, d := range m.Data {
			// for showcase purposes
			time.Sleep(100 * time.Millisecond)

			message := strings.Join([]string{strconv.FormatInt(d.Time, 10), d.Price}, ",")

			err := producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte(message),
			}, nil)
			if err != nil {
				fmt.Printf("Failed to produce message: %v\n", err)
			}
		}
	}

	producer.Flush(15 * 1000)
}
