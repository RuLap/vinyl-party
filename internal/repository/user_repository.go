package repository

import (
	"context"
	"vinyl-party/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id string) (*entity.User, error)
	GetByIDs(ids []string) ([]*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (r *userRepository) Create(user *entity.User) error {
	_, err := r.collection.InsertOne(context.Background(), user)
	return err
}

func (r *userRepository) GetByID(id string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByIDs(ids []string) ([]*entity.User, error) {
	var users []*entity.User
	cursor, err := r.collection.Find(context.Background(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user entity.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
