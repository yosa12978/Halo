package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func InitMongo(addr string, dbname string) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(addr)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	db = client.Database(dbname)
	return db, nil
}

func GetDB() *mongo.Database {
	return db
}
