package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-co-op/gocron"
	"os"
	"time"
)

const (
	topic             = "bitcoin"
	numPartitions     = 1
	replicationFactor = 1
)

var (
	bootstrapServers = os.Getenv("BOOTSTRAP_SERVERS")
)

func main() {
	fmt.Println(bootstrapServers)
	a := NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	defer a.Close()

	CreateTopics(a, []kafka.TopicSpecification{{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor}},
	)

	p := NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	defer p.Close()

	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Hour().Do(func() {
		data := GetData(topic)
		Produce(p, data, topic)
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	s.StartBlocking()
}
