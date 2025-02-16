package repository

import (
	"context"
	"vinyl-party/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByIDs(ctx context.Context, ids []string) ([]*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByIDs(ctx context.Context, ids []string) ([]*entity.User, error) {
	var users []*entity.User
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}
	
	return users, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
