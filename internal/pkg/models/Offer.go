package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Offer struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	Title        string             `bson:"title" json:"title"`
	Body         string             `bson:"body" json:"body"`
	Tags         []string           `bson:"tags" json:"tags"`
	Company      Company            `bson:"company" json:"company"`
	City         string             `bson:"city" json:"city"`
	ContactEmail string             `bson:"contactEmail" json:"contactEmail"`
	Website      string             `bson:"website" json:"website"`
}
