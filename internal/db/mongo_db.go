package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	client    *mongo.Client
	dbName    string
	validColl *mongo.Collection
	quarColl  *mongo.Collection
}

func NewMongoDB(client *mongo.Client, dbName string) *MongoDB {
	return &MongoDB{
		client:    client,
		dbName:    dbName,
		validColl: client.Database(dbName).Collection(CollectionValid),
		quarColl:  client.Database(dbName).Collection(CollectionQuarantine),
	}
}

func (m *MongoDB) InsertToValid(ctx context.Context, data interface{}, userID string) error {
	fmt.Println("InsertToValid in MongoDB")
	_, err := m.validColl.InsertOne(ctx, bson.M{
		"data":   data,
		"userID": userID,
	})
	return err
}

func (m *MongoDB) InsertToQuarantine(ctx context.Context, data interface{}, userID string, reason string) error {
	fmt.Println("InsertToQuarantine in MongoDB")
	_, err := m.quarColl.InsertOne(ctx, bson.M{
		"data":   data,
		"userID": userID,
		"reason": reason,
	})
	return err
}

func (m *MongoDB) GetAllValidData(ctx context.Context) ([]bson.M, error) {
	cursor, err := m.validColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
