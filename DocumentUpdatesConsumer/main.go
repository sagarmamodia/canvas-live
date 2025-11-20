package main

import (
	"DocumentUpdatesConsumer/config"
	"DocumentUpdatesConsumer/handler"
	"DocumentUpdatesConsumer/repository"
	"DocumentUpdatesConsumer/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	kafkaBroker = "localhost:9092"
	topic       = "document-updates"
	groupID     = "document-updates-consumer-group"
)

func connectDB(uri string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse the URI and setup client options
	clientOptions := options.Client().ApplyURI(uri)

	// connect
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	// ping the database to verify the connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}

	fmt.Println("Successfully connected to MongoDB!")
	return client
}

func main() {
	// Connect to DB
	client := connectDB(config.MongoConfig.MongoUri)

	// Repository
	r := repository.NewDocumentRepository(client, config.MongoConfig.DatabaseName, config.MongoConfig.DocumentCollectionName)

	// Create a new Kafka consumer
	fmt.Println("Trying to connect to Kafka!")
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		return
	}
	defer c.Close()
	fmt.Println("Connectd to Kafka!")

	// subscribe to a kafka topic
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
		return
	}

	// Setup a channel to handle OS signals for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	// Start consuming messages
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Received signal %v: terminating\n", sig)
			run = false
		default:
			// Poll for Kafka messages
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				// Process the consumed message
				fmt.Printf("Received message from topic %s: %s\n", *e.TopicPartition.Topic, string(e.Value))

				// Parse message into struct
				var msg types.Message
				if err := json.Unmarshal(e.Value, &msg); err != nil {
					fmt.Printf("[Error] Can't unmarshall message")
				}

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

				go func() {
					defer cancel()
					handler.DocumentUpdatesHandler(ctx, r, msg)
				}()

			case kafka.Error:
				// Handle Kafka errors
				fmt.Printf("Error: %v\n", e)
			}
		}
	}
}
