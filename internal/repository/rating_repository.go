package repository

import (
	"context"
	"vinyl-party/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatingRepository interface {
	Create(rating *entity.Rating) error
	GetByID(id string) (*entity.Rating, error)
	GetByIDs(ids []string) ([]*entity.Rating, error)
}

type ratingRepository struct {
	collection *mongo.Collection
}

func NewRatingRepository(db *mongo.Database) RatingRepository {
	return &ratingRepository{
		collection: db.Collection("ratings"),
	}
}

func (r *ratingRepository) Create(rating *entity.Rating) error {
	_, err := r.collection.InsertOne(context.Background(), rating)
	return err
}

func (r *ratingRepository) GetByID(id string) (*entity.Rating, error) {
	var rating entity.Rating
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&rating)
	if err != nil {
		return nil, err
	}

	return &rating, nil
}

func (r *ratingRepository) GetByIDs(ids []string) ([]*entity.Rating, error) {
	var ratings []*entity.Rating
	cursor, err := r.collection.Find(context.Background(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &ratings); err != nil {
		return nil, err
	}

	return ratings, nil
}
