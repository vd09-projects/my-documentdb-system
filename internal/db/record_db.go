package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/vd09-projects/my-documentdb-system/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecordDB struct {
	client       *mongo.Client
	dbName       string
	validColl    *mongo.Collection
	quarColl     *mongo.Collection
	recordFields *mongo.Collection
}

func NewRecordDB(client *mongo.Client, dbName string) *RecordDB {
	return &RecordDB{
		client:       client,
		dbName:       dbName,
		validColl:    client.Database(dbName).Collection(ValidCollection),
		quarColl:     client.Database(dbName).Collection(QuarantineCollection),
		recordFields: client.Database(dbName).Collection(RecordFieldsCollection),
	}
}

func (m *RecordDB) InsertToValid(ctx context.Context, data interface{}, userID string, recordType string) error {
	_, err := m.validColl.InsertOne(ctx, bson.M{
		"data":       data,
		"userID":     userID,
		"recordType": recordType,
		"timestamp":  time.Now(),
	})
	if err != nil {
		return err
	}

	return m.InsertToRecordFields(ctx, data, userID, recordType)
}

func (m *RecordDB) InsertToRecordFields(ctx context.Context, data interface{}, userID string, recordType string) error {
	// Step 1: Define the filter to find the document with the given userID and recordType
	filter := bson.M{
		"userID":     userID,
		"recordType": recordType,
	}

	// Step 2: Define the update to merge the new fields into the existing set
	update := bson.M{
		"$addToSet": bson.M{
			"fields": bson.M{
				"$each": utils.TraverseDynamicJSON(data),
			},
		},
	}

	// Step 4: Specify the upsert option
	options := options.Update().SetUpsert(true)

	// Step 5: Perform the upsert operation
	_, err := m.recordFields.UpdateOne(ctx, filter, update, options)
	if err != nil {
		return fmt.Errorf("failed to upsert fields for userID %s and recordType %s: %w", userID, recordType, err)
	}

	return nil
}

func (m *RecordDB) InsertToQuarantine(ctx context.Context, data interface{}, userID string, recordType string, reason string) error {
	_, err := m.quarColl.InsertOne(ctx, bson.M{
		"data":       data,
		"userID":     userID,
		"recordType": recordType,
		"reason":     reason,
		"timestamp":  time.Now(),
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
		// Stage 2: Project fields to include recordType and data
		{{Key: "$project", Value: bson.M{
			"recordType": 1,
			"data":       1,
			"_id":        0,
		}}},
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

// GetRecordTypesForUser fetches distinct record types for a user.
func (m *RecordDB) GetRecordTypesForUser(ctx context.Context, userID string) ([]bson.M, error) {
	// Build the base filter for user ID
	filter := bson.M{
		"userID": userID,
	}

	// Create the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Stage 1: Match the filter
		{{Key: "$match", Value: filter}},
		// Stage 2: Group by recordType to get distinct values
		{{Key: "$group", Value: bson.M{
			"_id": "$recordType", // Group by recordType
		}}},
		// Stage 3: Project the result to include recordType field
		{{Key: "$project", Value: bson.M{
			"recordType": "$_id",
			"_id":        0,
		}}},
	}

	// Execute the aggregation pipeline
	cursor, err := m.recordFields.Aggregate(ctx, pipeline)
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

// GetFieldsForUserAndType fetches distinct firlds by record types and user.
func (m *RecordDB) GetFieldsForUserAndType(ctx context.Context, userID string, recordType string) (bson.M, error) {
	// Build the base filter for user ID
	filter := bson.M{
		"userID":     userID,
		"recordType": recordType,
	}

	// Create the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Stage 1: Match the filter
		{{Key: "$match", Value: filter}},
		// Stage 3: Project the result to include fields field
		{{Key: "$project", Value: bson.M{
			"fields": 1,
			"_id":    0,
		}}},
	}

	// Execute the aggregation pipeline
	cursor, err := m.recordFields.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results[0], nil
}

func (m *RecordDB) AggregateData(ctx context.Context, userID, recordType, field string) ([]float64, error) {
	// Build the base filter
	filter := bson.M{
		"userID":     userID,
		"recordType": recordType,
	}

	// Fetch raw data from MongoDB
	cursor, err := m.validColl.Find(ctx, filter)
	if err != nil {
		return []float64{}, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer cursor.Close(ctx)

	// Parse the results into a slice of bson.M
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return []float64{}, fmt.Errorf("failed to read data: %w", err)
	}

	// Apply aggregation in code
	var values []float64
	for _, record := range results {
		// Extract the field value
		rawValue, ok := record["data"].(bson.M)[field]
		if !ok {
			continue // Skip if the field is not found
		}

		// Convert the value to float64 if possible
		switch v := rawValue.(type) {
		case string:
			floatValue, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return []float64{}, fmt.Errorf("failed to convert field value to float: %w", err)
			}
			values = append(values, floatValue)
		case float64:
			values = append(values, v)
		default:
			return []float64{}, fmt.Errorf("unsupported field value type: %T", v)
		}
	}
	return values, nil
}
