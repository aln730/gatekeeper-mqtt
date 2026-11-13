package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func checkAccess(userID, doorID string) (*bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user bson.M
	err := db.Collection("users").FindOne(ctx, bson.M{"id": userID}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if user["disabled"] == true {
		granted := false
		return &granted, nil
	}

	var ticket bson.M
	err = db.Collection("userTickets").FindOne(ctx, bson.M{
		"userId": bson.M{"$in": bson.A{userID, "*"}},
		"doorId": bson.M{"$in": bson.A{doorID, "*"}},
	}, options.FindOne().SetSort(bson.D{{Key: "priority", Value: -1}})).Decode(&ticket)

	if err == nil {
		granted := ticket["granted"].(bool)
		return &granted, nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	groups, _ := user["groups"].(bson.A)

	err = db.Collection("groupTickets").FindOne(ctx, bson.M{
		"doorId":  bson.M{"$in": bson.A{"*", doorID}},
		"groupId": bson.M{"$in": append(groups, "*")},
	}, options.FindOne().SetSort(bson.D{{Key: "priority", Value: -1}})).Decode(&ticket)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	granted := ticket["granted"].(bool)
	return &granted, nil
}
