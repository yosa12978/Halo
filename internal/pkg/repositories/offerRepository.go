package repositories

import (
	"context"
	"errors"
	"sync"

	"github.com/yosa12978/halo/internal/pkg/dto"
	"github.com/yosa12978/halo/internal/pkg/models"
	mongodb "github.com/yosa12978/halo/internal/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IOfferRepository interface {
	CreateOffer(user_id string, uo dto.CreateOffer) error
	GetOffer(id_hex string) (*models.Offer, error)
	SearchOffers(q string) []models.Offer
	UpdateOffer(id_hex string, user_id string, uo dto.UpdateOffer) error
	DeleteOffer(id_hex string, user_id string) error
	ActiveOffer(id_hex string, user_id_hex string) error
}

type OfferRepository struct {
	Db *mongo.Database
}

func NewOfferRepository() IOfferRepository {
	return &OfferRepository{Db: mongodb.GetDB()}
}

func (or *OfferRepository) GetOffer(id_hex string) (*models.Offer, error) {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return nil, errors.New("offer not found")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var offer models.Offer
	err = or.Db.Collection("offers").FindOne(ctx, bson.M{"_id": id}).Decode(&offer)
	if err != nil {
		return nil, errors.New("offer not found")
	}
	return &offer, nil
}

func (or *OfferRepository) SearchOffers(q string) []models.Offer {
	opt := options.Find().SetSort(bson.M{"_id": -1})
	filter := bson.M{
		"title":               bson.M{"$regex": primitive.Regex{Pattern: q}},
		"body":                bson.M{"$regex": primitive.Regex{Pattern: q}},
		"tags":                bson.M{"$regex": primitive.Regex{Pattern: q}},
		"company.companyName": bson.M{"$regex": primitive.Regex{Pattern: q}},
		"city":                bson.M{"$regex": primitive.Regex{Pattern: q}},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var offers []models.Offer
	cursor, err := or.Db.Collection("offers").Find(ctx, filter, opt)
	if err != nil {
		return offers
	}

	var wg sync.WaitGroup
	for cursor.Next(ctx) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var offer models.Offer
			cursor.Decode(&offer)
			offers = append(offers, offer)
		}()
	}
	wg.Wait()

	return offers
}

func (or *OfferRepository) CreateOffer(user_id string, uo dto.CreateOffer) error {
	companyRepository := NewCompanyRepository()
	comp, err := companyRepository.GetCompany(uo.Company)
	offer := models.Offer{
		ID:           primitive.NewObjectID(),
		Title:        uo.Title,
		Body:         uo.Body,
		Tags:         uo.Tags,
		Company:      *comp,
		City:         uo.City,
		ContactEmail: uo.ContactEmail,
		Website:      uo.Website,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = or.Db.Collection("offers").InsertOne(ctx, offer)
	return err
}

func (or *OfferRepository) UpdateOffer(id_hex string, user_id_hex string, uo dto.UpdateOffer) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return err
	}
	user_id, err := primitive.ObjectIDFromHex(user_id_hex)
	if err != nil {
		return err
	}
	upd := bson.M{"$set": []bson.M{
		{"title": uo.Title},
		{"body": uo.Body},
		{"tags": uo.Tags},
		{"city": uo.City},
		{"contactEmail": uo.ContactEmail},
	},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err = or.Db.Collection("offers").UpdateOne(ctx, bson.M{"_id": id, "company._id": user_id}, upd)
	return err
}

func (or *OfferRepository) DeleteOffer(id_hex string, user_id_hex string) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return err
	}
	user_id, err := primitive.ObjectIDFromHex(user_id_hex)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	filter := bson.M{"$and": []bson.M{{"_id": id}, {"company._id": user_id}}}
	_, err = or.Db.Collection("offers").DeleteOne(ctx, filter)
	return err
}

func (or *OfferRepository) ActiveOffer(id_hex string, user_id_hex string) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return errors.New("offer not found")
	}
	user_id, err := primitive.ObjectIDFromHex(user_id_hex)
	if err != nil {
		return errors.New("user not found")
	}
	offer, err := or.GetOffer(id_hex)
	if err != nil {
		return err
	}
	upd := bson.M{"$set": []bson.M{{"active": !offer.Active}}}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	filter := bson.M{"$and": []bson.M{{"_id": id}, {"company._id": user_id}}}
	_, err = or.Db.Collection("offers").UpdateOne(ctx, filter, upd)
	return err
}
