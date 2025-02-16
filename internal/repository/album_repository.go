package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"vinyl-party/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlbumRepository interface {
	EnsureIndexes() error
	Create(ctx context.Context, album *entity.Album) error
	GetByID(ctx context.Context, id string) (*entity.Album, error)
	GetByIDs(ctx context.Context, ids []string) ([]*entity.Album, error)
	AddRating(ctx context.Context, albumID string, ratingID string) error
	Delete(ctx context.Context, id string) error
}

type albumRepository struct {
	collection *mongo.Collection
}

func NewAlbumRepository(db *mongo.Database) AlbumRepository {
	return &albumRepository{
		collection: db.Collection("albums"),
	}
}

func (r *albumRepository) EnsureIndexes() error {
	ratingIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "participants.user_id", Value: 1},
			{Key: "date", Value: 1},
		},
		Options: options.Index().SetName("user_parties"),
	}

	uniqueRating := mongo.IndexModel{
		Keys: bson.D{
			{Key: "_id", Value: 1},
			{Key: "participants.user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	indexes := []mongo.IndexModel{
		ratingIndex,
		uniqueRating,
	}
	_, err := r.collection.Indexes().CreateMany(context.Background(), indexes)

	return err
}

func (r *albumRepository) Create(ctx context.Context, album *entity.Album) error {
	_, err := r.collection.InsertOne(ctx, album)
	return err
}

func (r *albumRepository) GetByID(ctx context.Context, id string) (*entity.Album, error) {
	var album entity.Album
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&album)
	if err != nil {
		return nil, err
	}

	return &album, nil
}

func (r *albumRepository) GetByIDs(ctx context.Context, ids []string) ([]*entity.Album, error) {
	var albums []*entity.Album
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &albums)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (r *albumRepository) AddRating(ctx context.Context, albumID string, ratingID string) error {
	update := bson.M{"$addToSet": bson.M{"rating_ids": ratingID}}
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": albumID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("album not found")
	}

	return nil
}

func (r *albumRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("album not found")
	}

	return nil
}
