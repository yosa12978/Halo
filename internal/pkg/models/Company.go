package models

type Company struct {
	BaseAccount
	CompanyName    string `bson:"companyName" json:"companyName"`
	CompanyWebsite string `bson:"website" json:"website"`
	ContactEmail   string `bson:"contactEmail" json:"contactEmail"`
}
