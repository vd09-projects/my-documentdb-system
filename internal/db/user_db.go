package db

import (
	"context"
	"errors"

	"github.com/vd09-projects/my-documentdb-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const UsersCollection = "users"

type UserDB struct {
	coll *mongo.Collection
}

func NewUserDB(client *mongo.Client, dbName string) *UserDB {
	return &UserDB{
		coll: client.Database(dbName).Collection(UsersCollection),
	}
}

// CreateUser inserts a new user record (with a hashed password)
func (u *UserDB) CreateUser(ctx context.Context, user models.User) error {
	// Check if username already exists
	filter := bson.M{"username": user.Username}
	err := u.coll.FindOne(ctx, filter).Err()
	if err == nil {
		return errors.New("username already taken")
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	_, err = u.coll.InsertOne(ctx, user)
	return err
}

// FindUserByUsername returns the User record for the given username
func (u *UserDB) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	filter := bson.M{"username": username}
	var user models.User
	err := u.coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
