package repositories

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/yosa12978/halo/internal/pkg/dto"
	"github.com/yosa12978/halo/internal/pkg/models"
	"github.com/yosa12978/halo/pkg/crypto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserRepository interface {
	Login(lu dto.LoginUser) (string, error)
	CreateUser(cu dto.CreateUser) error
	GetUsers() []models.User
	GetUserByID(id string) *models.User
	GetUserByEmailAndPwd(email string, pwd string) *models.User
	GetUserByEmail(email string) *models.User
	UpdateUser(id_hex string, pwd string, uu dto.UpdateUser) error
	DeleteUser(id string, pwd string) error
}

type UserRepository struct {
	Db *mongo.Database
}

func NewUserRepository(db *mongo.Database) IUserRepository {
	return &UserRepository{Db: db}
}

func (ur *UserRepository) Login(lu dto.LoginUser) (string, error) {
	return "", nil
}

func (ur *UserRepository) CreateUser(cu dto.CreateUser) error {
	usr := ur.GetUserByEmail(cu.Email)
	compRepo := NewCompanyRepository(ur.Db)
	co := compRepo.GetCompanyByEmail(cu.Email)
	if co != nil || usr != nil {
		return errors.New("email is already in use")
	}
	user := models.User{
		BaseAccount: models.BaseAccount{
			ID:       primitive.NewObjectID(),
			Email:    cu.Email,
			Password: crypto.GetMD5(cu.Password),
			Roles:    []string{models.ROLE_USER},
			Regdate:  time.Now().Unix(),
		},
		FirstName: cu.FirstName,
		LastName:  cu.LastName,
		BirthDate: cu.BirthDate,
		City:      cu.City,
		Country:   cu.Country,
		Skills:    cu.Skills,
		Status:    models.LOOKING,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ur.Db.Collection("users").InsertOne(ctx, user)
	return nil
}

func (ur *UserRepository) GetUsers() []models.User {
	findOptions := options.Find().SetSort(bson.M{"_id": -1})
	var users []models.User
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cursor, err := ur.Db.Collection("users").Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return users
	}

	var wg sync.WaitGroup
	for cursor.Next(ctx) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var user models.User
			cursor.Decode(&user)
			users = append(users, user)
		}()
	}
	wg.Wait()

	return users
}

func (ur *UserRepository) GetUserByID(id_hex string) *models.User {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var user models.User
	err = ur.Db.Collection("users").FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil
	}
	return &user
}

func (ur *UserRepository) GetUserByEmail(email string) *models.User {
	var user models.User
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := ur.Db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil
	}
	return &user
}

func (ur *UserRepository) GetUserByEmailAndPwd(email string, pwd string) *models.User {
	var user models.User
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	filter := bson.M{"$and": []bson.M{{"email": email}, {"password": crypto.GetMD5(pwd)}}}
	err := ur.Db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil
	}
	return &user
}

func (ur *UserRepository) UpdateUser(id_hex string, pwd string, uu dto.UpdateUser) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	upd := bson.M{"$set": []bson.M{
		{"firstName": uu.FirstName},
		{"lastName": uu.LastName},
		{"birthDate": uu.BirthDate},
		{"city": uu.City},
		{"country": uu.Country},
		{"skills": uu.Skills},
		{"statis": uu.Status},
	},
	}
	_, err = ur.Db.Collection("users").UpdateOne(ctx, bson.M{"_id": id}, upd)
	return err
}

func (ur *UserRepository) DeleteUser(id_hex string, pwd string) error {
	id, err := primitive.ObjectIDFromHex(id_hex)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	filter := bson.M{"$and": []bson.M{{"_id": id}, {"password": pwd}}}
	_, err = ur.Db.Collection("users").DeleteOne(ctx, filter)
	return err
}
