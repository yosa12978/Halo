package repositories

import "go.mongodb.org/mongo-driver/mongo"

type OfferRepository struct {
	Db *mongo.Database
}
