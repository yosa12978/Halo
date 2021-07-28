package repositories

import "go.mongodb.org/mongo-driver/mongo"

type UserRepository struct {
	Db *mongo.Database
}
