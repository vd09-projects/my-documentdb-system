package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Database interface {
	InsertToValid(ctx context.Context, data interface{}, userID string, recordType string) error
	InsertToRecordFields(ctx context.Context, data interface{}, userID string, recordType string) error
	InsertToQuarantine(ctx context.Context, data interface{}, userID string, recordType string, reason string) error
	GetAllValidData(ctx context.Context) ([]bson.M, error)
	GetUserData(ctx context.Context, userID string, from, to *time.Time) ([]bson.M, error)
	GetRecordTypesForUser(ctx context.Context, userID string) ([]bson.M, error)
	GetFieldsForUserAndType(ctx context.Context, userID string, recordType string) (bson.M, error)
	AggregateData(ctx context.Context, userID, recordType, field, op string) (float64, error)
}
