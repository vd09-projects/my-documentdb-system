package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecordDB struct {
	client    *mongo.Client
	dbName    string
	validColl *mongo.Collection
	quarColl  *mongo.Collection
}

func NewRecordDB(client *mongo.Client, dbName string) *RecordDB {
	return &RecordDB{
		client:    client,
		dbName:    dbName,
		validColl: client.Database(dbName).Collection(ValidCollection),
		quarColl:  client.Database(dbName).Collection(QuarantineCollection),
	}
}

func (m *RecordDB) InsertToValid(ctx context.Context, data interface{}, userID string) error {
	fmt.Println("InsertToValid in MongoDB")
	_, err := m.validColl.InsertOne(ctx, bson.M{
		"data":      data,
		"userID":    userID,
		"timestamp": time.Now(),
	})
	return err
}

func (m *RecordDB) InsertToQuarantine(ctx context.Context, data interface{}, userID string, reason string) error {
	fmt.Println("InsertToQuarantine in MongoDB")
	_, err := m.quarColl.InsertOne(ctx, bson.M{
		"data":      data,
		"userID":    userID,
		"reason":    reason,
		"timestamp": time.Now(),
	})
	return err
}

func (m *RecordDB) GetAllValidData(ctx context.Context) ([]bson.M, error) {
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

// GetUserData fetches user-specific data, optionally filtered by date range.
func (m *RecordDB) GetUserData(ctx context.Context, userID string, from, to *time.Time) ([]bson.M, error) {
	// Build the base filter for user ID
	filter := bson.M{
		"userID": userID,
	}

	// If you’re using a 'timestamp' field, add date range conditions
	// only if 'from' or 'to' are provided
	dateFilter := bson.M{}
	if from != nil {
		dateFilter["$gte"] = *from
	}
	if to != nil {
		dateFilter["$lte"] = *to
	}
	if len(dateFilter) > 0 {
		filter["timestamp"] = dateFilter
	}

	// Create the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Stage 1: Match the filter
		{{Key: "$match", Value: filter}},
		// Stage 2: Replace the root with the 'data' field
		{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$data"}}},
	}

	// Execute the aggregation pipeline
	cursor, err := m.validColl.Aggregate(ctx, pipeline)
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
