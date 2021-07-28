package models

import (
	"time"
)

type User struct {
	BaseAccount
	FirstName string    `bson:"firstName" json:"firstName"`
	LastName  string    `bson:"lastName" json:"lastName"`
	BirthDate time.Time `bson:"birthDate" json:"birthDate"`
	City      string    `bson:"city" json:"city"`
	Country   string    `bson:"country" json:"country"`
}
