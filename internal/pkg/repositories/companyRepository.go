package repositories

import (
	"context"
	"sync"
	"time"

	"github.com/yosa12978/halo/internal/pkg/dto"
	"github.com/yosa12978/halo/internal/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICompanyRepository interface {
	GetCompanies() []models.Company
	GetCompany(id_hex string) *models.Company
	SearchCompany(s string) []models.Company
	CreateCompany(comp dto.CreateCompany) error
	UpdateCompany(id_hex string, comp dto.UpdateCompany) error
	DeleteCompany(id_hex string) error
}

type CompanyRepository struct {
	Db *mongodb.Database
}

func NewCompanyRepository(db *mongodb.Database) ICompanyRepository {
	return &CompanyRepository{Db: db}
}

func (cr *CompanyRepository) GetCompanies() []models.Company {
	findOptions := options.Find().SetSort(bson.M{"_id": -1})
	var companies []models.Company
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cursor, err := cr.Db.Collection("companies").Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return companies
	}

	var wg sync.WaitGroup
	for cursor.Next(ctx) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var company models.Company
			cursor.Decode(&company)
			companies = append(companies, company)
		}()
	}
	wg.Wait()

	return companies
}

func (cr *CompanyRepository) GetCompany(id_hex string) *models.Company {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return nil
	}
	var company models.Company
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = cr.Db.Collection("companies").FindOne(ctx, bson.M{"_id": id}).Decode(&company)
	if err != nil {
		return nil
	}
	return &company
}

func (cr *CompanyRepository) SearchCompany(s string) []models.Company {
	var companies []models.Company
	filter := bson.M{"title": bson.M{"$regex": primitive.Regex{Pattern: s}}}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cursor, err := cr.Db.Collection("companies").Find(ctx, filter)
	if err != nil {
		return companies
	}

	var wg sync.WaitGroup
	for cursor.Next(ctx) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var company models.Company
			cursor.Decode(&company)
			companies = append(companies, company)
		}()
	}
	wg.Wait()

	return companies
}

func (cr *CompanyRepository) CreateCompany(comp dto.CreateCompany) error {
	company := models.Company{
		BaseAccount: models.BaseAccount{
			ID:       primitive.NewObjectID(),
			Email:    comp.Email,
			Password: comp.Password,
			Roles:    []string{models.ROLE_COMPANY},
			Regdate:  time.Now().Unix(),
		},
		CompanyName:    comp.Name,
		CompanyWebsite: comp.Website,
		ContactEmail:   comp.Contact,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := cr.Db.Collection("companies").InsertOne(ctx, company)
	if err != nil {
		return err
	}
	return nil
}

func (cr *CompanyRepository) UpdateCompany(id_hex string, comp dto.UpdateCompany) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	upd := bson.M{"$set": []bson.M{
		{"companyName": comp.Name},
		{"contactEmail": comp.Email},
		{"website": comp.Website},
		{"email": comp.Email},
	},
	}
	_, err = cr.Db.Collection("companies").UpdateOne(ctx, bson.M{"_id": id}, upd)
	return err
}

func (cr *CompanyRepository) DeleteCompany(id_hex string) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = cr.Db.Collection("companies").DeleteOne(ctx, bson.M{"_id": id})
	return err
}
