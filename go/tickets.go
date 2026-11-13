
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func getTickets(w http.ResponseWriter, r *http.Request) {
	var collection string

	switch r.PathValue("variant") {
	case "groups":
		collection = "groupTickets"
	case "users":
		collection = "userTickets"
	default:
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{
		{Key: "priority", Value: -1},
		{Key: "_id", Value: 1},
	})

	cursor, err := db.Collection(collection).Find(ctx, bson.D{}, opts)
	if err != nil {
		http.Error(w, "Failed to fetch tickets", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	tickets := []bson.M{}
	if err := cursor.All(ctx, &tickets); err != nil {
		http.Error(w, "Failed to read tickets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickets)
}