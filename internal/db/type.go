package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Database interface {
	InsertToValid(ctx context.Context, data interface{}, userID string) error
	InsertToQuarantine(ctx context.Context, data interface{}, userID string, reason string) error
	GetAllValidData(ctx context.Context) ([]bson.M, error)
	GetUserData(ctx context.Context, userID string, from, to *time.Time) ([]bson.M, error)
}
