package main

import (
	"UpdatesReceiverAndTransmitterService/kafkaUtils"
	"UpdatesReceiverAndTransmitterService/redis"
	"UpdatesReceiverAndTransmitterService/websocket"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

func wsHandler(pool *websocket.Pool, redis_client *redis.RedisClient) gin.HandlerFunc {
	// Return a Gin handler function
	return func(c *gin.Context) {
		docId := c.Param("docId")
		if docId == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "documentId missing"})
		}
		// 1. Authentication Check (Using c.Request)
		// Access header directly from the raw http.Request object
		userId := c.Request.Header.Get("X-User-ID")
		username := c.Request.Header.Get("X-Username")
		if userId == "" {
			// Use Gin's method to send HTTP error response before upgrade
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			return
		}

		// 2. Perform WebSocket Upgrade (Using c.Writer and c.Request)
		conn, err := websocket.Upgrade(c.Writer, c.Request)
		if err != nil {
			// Log error after upgrade attempt, as headers may already be sent
			log.Printf("WebSocket Upgrade Failed: %v", err)
			// Note: Since upgrade failed, you cannot use c.JSON here
			return
		}

		// 3. Initialize and Register Client
		client := &websocket.Client{
			UserID:      userId,
			Username:    username,
			DocumentID:  docId, // Ensure this is correctly retrieved or set
			Conn:        conn,
			Pool:        pool,
			Send:        make(chan []byte),
			RedisClient: redis_client,
		}

		fmt.Println("[WsHandler] client reader running!")
		go client.Writer() // Start a goroutine responsible for send message(it receives via Send channel) to the client
		fmt.Println("[WsHandler] client Writer running!")

		pool.Register <- client
		client.Read() // Start the client's read loop
	}
}

func main() {
	// kafka Setup
	fmt.Println("Trying to connect to Kafka!")
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaUtils.KafkaBroker})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return
	}
	defer p.Close()
	fmt.Println("Connected to Kafka!")

	// Redis Setup
	redis_client := redis.NewRedisClient("localhost:6379")

	// Websocket pool
	pool := websocket.NewPool(p)
	go pool.Start()

	// Server setup
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server running.")
	})

	router.GET("/ws/:docId", wsHandler(pool, redis_client))

	router.Run(":8083")
}
