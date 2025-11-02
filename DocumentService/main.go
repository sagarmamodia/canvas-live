package main

import (
	"context"
	"document-service/config"
	"document-service/handler"
	"document-service/middleware"
	"document-service/repository"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// Set up Repositories
	DocumentRepository := repository.NewDocumentRepository(client, config.MongoConfig.DatabaseName, config.MongoConfig.DocumentCollectionName, config.MongoConfig.SharedDocRecordCollectionName)

	// Set up Handlers
	documentHandler := handler.DocumentHandler{DocumentRepository: DocumentRepository}

	// Server
	mux := http.NewServeMux()
	mux.Handle("/document/create", http.HandlerFunc(documentHandler.CreateNewDocument))
	mux.Handle("/document/all", http.HandlerFunc(documentHandler.GetAllDocuments))
	// mux.Handle("/document/share", http.HandlerFunc(documentHandler.CreateNewDocument))
	// mux.Handle("/document/create", http.HandlerFunc(documentHandler.CreateNewDocument))
	// mux.Handle("/document/create", http.HandlerFunc(documentHandler.CreateNewDocument))

	finalMux := middleware.RequestLoggingMiddleware(mux)
	fmt.Println("Starting server on port 8081")

	if err := http.ListenAndServe(":8081", finalMux); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}
