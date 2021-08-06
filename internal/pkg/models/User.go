package models

type User struct {
	BaseAccount
	FirstName string   `bson:"firstName" json:"firstName"`
	LastName  string   `bson:"lastName" json:"lastName"`
	BirthDate int64    `bson:"birthDate" json:"birthDate"`
	City      string   `bson:"city" json:"city"`
	Country   string   `bson:"country" json:"country"`
	Skills    []string `bson:"skills" json:"skills"`
	Status    int8     `bson:"status" json:"status"`
}
