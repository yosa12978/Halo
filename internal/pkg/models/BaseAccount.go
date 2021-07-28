package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BaseAccount struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
	Roles    []string           `bson:"roles" json:"roles"`
	Regdate  int64              `bson:"regdate" json:"regdate"`
}
