package repositories

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/yosa12978/halo/internal/pkg/dto"
	"github.com/yosa12978/halo/internal/pkg/models"
	mongodb "github.com/yosa12978/halo/internal/pkg/mongo"
	"github.com/yosa12978/halo/pkg/crypto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICompanyRepository interface {
	GetCompanies() []models.Company
	GetCompany(id_hex string) (*models.Company, error)
	GetCompanyByEmail(email string) (*models.Company, error)
	GetCompanyByEmailAndPass(email string, password string) (*models.Company, error)
	SearchCompany(s string) []models.Company
	CreateCompany(comp dto.CreateCompany) error
	UpdateCompany(id_hex string, comp dto.UpdateCompany) error
	DeleteCompany(id_hex string) error
	IsCompanyExist(email string) bool
	IsCompPassExist(email string, password string) bool
}

type CompanyRepository struct {
	Db *mongo.Database
}

func NewCompanyRepository() ICompanyRepository {
	return &CompanyRepository{Db: mongodb.GetDB()}
}

func (cr *CompanyRepository) GetCompanies() []models.Company {
	findOptions := options.Find().SetSort(bson.M{"baseaccount._id": -1})
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

func (cr *CompanyRepository) GetCompany(id_hex string) (*models.Company, error) {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return nil, errors.New("company not found")
	}
	var company models.Company
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = cr.Db.Collection("companies").FindOne(ctx, bson.M{"baseaccount._id": id}).Decode(&company)
	if err != nil {
		return nil, errors.New("company not found")
	}
	return &company, nil
}

func (cr *CompanyRepository) GetCompanyByEmail(email string) (*models.Company, error) {
	var company models.Company
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := cr.Db.Collection("companies").FindOne(ctx, bson.M{"baseaccount.email": email}).Decode(&company)
	if err != nil {
		return nil, errors.New("company not found")
	}
	return &company, nil
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
	co := cr.IsCompanyExist(comp.Email)
	usrRepo := NewUserRepository()
	usr := usrRepo.IsUserExist(comp.Email)
	if co || usr {
		return errors.New("email is already in use")
	}

	company := models.Company{
		BaseAccount: models.BaseAccount{
			ID:       primitive.NewObjectID(),
			Email:    comp.Email,
			Password: crypto.GetMD5(comp.Password),
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
		{"contactEmail": comp.Contact},
		{"website": comp.Website},
		{"email": comp.Email},
	},
	}
	_, err = cr.Db.Collection("companies").UpdateOne(ctx, bson.M{"baseaccount._id": id}, upd)
	return err
}

func (cr *CompanyRepository) DeleteCompany(id_hex string) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = cr.Db.Collection("companies").DeleteOne(ctx, bson.M{"baseaccount._id": id})
	return err
}

func (cr *CompanyRepository) IsCompanyExist(email string) bool {
	var company models.Company
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := cr.Db.Collection("companies").FindOne(ctx, bson.M{"baseaccount.email": email}).Decode(&company)
	if err != nil {
		return false
	}
	return true
}

func (cr *CompanyRepository) GetCompanyByEmailAndPass(email string, password string) (*models.Company, error) {
	var company models.Company
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	filter := bson.M{
		"$and": []bson.M{
			{"baseaccount.email": email},
			{"baseaccount.password": crypto.GetMD5(password)},
		},
	}
	err := cr.Db.Collection("companies").FindOne(ctx, filter).Decode(&company)
	if err != nil {
		return nil, errors.New("company not found")
	}
	return &company, nil
}

func (cr *CompanyRepository) IsCompPassExist(email string, password string) bool {
	var comp models.Company
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	filter := bson.M{"$and": []bson.M{
		{"baseaccount.email": email},
		{"baseaccount.password": crypto.GetMD5(password)},
	},
	}
	err := cr.Db.Collection("companies").FindOne(ctx, filter).Decode(&comp)
	if err != nil {
		return false
	}
	return true
}
