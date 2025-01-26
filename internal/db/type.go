package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type Database interface {
	InsertToValid(ctx context.Context, data interface{}, userID string) error
	InsertToQuarantine(ctx context.Context, data interface{}, userID string, reason string) error
	GetAllValidData(ctx context.Context) ([]bson.M, error)
}
