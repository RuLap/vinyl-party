package repository

import (
	"context"
	"vinyl-party/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatingRepository interface {
	Create(ctx context.Context, rating *entity.Rating) error
	GetByID(ctx context.Context, id string) (*entity.Rating, error)
	GetByIDs(ctx context.Context, ids []string) ([]*entity.Rating, error)
	GetByAlbumID(ctx context.Context, albumID string) ([]*entity.Rating, error)
}

type ratingRepository struct {
	collection *mongo.Collection
}

func NewRatingRepository(db *mongo.Database) RatingRepository {
	return &ratingRepository{
		collection: db.Collection("ratings"),
	}
}

func (r *ratingRepository) Create(ctx context.Context, rating *entity.Rating) error {
	_, err := r.collection.InsertOne(ctx, rating)
	return err
}

func (r *ratingRepository) GetByID(ctx context.Context, id string) (*entity.Rating, error) {
	var rating entity.Rating
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&rating)
	if err != nil {
		return nil, err
	}

	return &rating, nil
}

func (r *ratingRepository) GetByIDs(ctx context.Context, ids []string) ([]*entity.Rating, error) {
	var ratings []*entity.Rating
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &ratings)
	if err != nil {
		return nil, err
	}

	return ratings, err
}

func (r *ratingRepository) GetByAlbumID(ctx context.Context, albumID string) ([]*entity.Rating, error) {
	var ratings []*entity.Rating
	cursor, err := r.collection.Find(ctx, bson.M{"album_id": albumID})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &ratings)
	if err != nil {
		return nil, err
	}

	return ratings, nil
}
