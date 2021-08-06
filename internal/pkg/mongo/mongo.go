package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func InitMongo(addr string, dbname string) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(addr)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("MongoDb connection attempt failed\n")
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	db = client.Database(dbname)
	log.Println("Connected to MongoDb")
	return db, nil
}

func GetDB() *mongo.Database {
	return db
}
