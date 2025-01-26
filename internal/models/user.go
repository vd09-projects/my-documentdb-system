package models

type User struct {
	ID       string `bson:"_id,omitempty"` // or use primitive.ObjectID if you prefer
	Username string `bson:"username"`
	Password string `bson:"password"` // hashed password
}
