
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var db *mongo.Database

func main() {
	uri := os.Getenv("GK_MONGO_SERVER")
	if uri == "" {
		log.Fatal("database env not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	db = client.Database("gatekeeper")

	log.Println("DB connection opened!")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tickets/{variant}", getTickets)

	log.Println("HTTP server listening on :3001")
	log.Fatal(http.ListenAndServe(":3001", mux))
}