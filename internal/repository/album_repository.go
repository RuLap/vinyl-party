package repository

import (
	"context"
	"errors"
	"vinyl-party/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlbumRepository interface {
	Create(album *entity.Album) error
	GetByID(id string) (*entity.Album, error)
	GetByPartyID(partyID string) ([]*entity.Album, error)
	AddRating(albumID string, ratingID string) error
	Delete(id string) error
}

type albumRepository struct {
	collection *mongo.Collection
}

func NewAlbumRepository(db *mongo.Database) AlbumRepository {
	return &albumRepository{
		collection: db.Collection("albums"),
	}
}

func (r *albumRepository) Create(album *entity.Album) error {
	_, err := r.collection.InsertOne(context.Background(), album)
	return err
}

func (r *albumRepository) GetByID(id string) (*entity.Album, error) {
	var album entity.Album
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&album)
	if err != nil {
		return nil, err
	}

	return &album, nil
}

func (r *albumRepository) GetByPartyID(partyID string) ([]*entity.Album, error) {
	var albums []*entity.Album
	cursor, err := r.collection.Find(context.Background(), bson.M{"party_id": partyID})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &albums)
	return albums, err
}

func (r *albumRepository) AddRating(albumID string, ratingID string) error {
	update := bson.M{"$addToSet": bson.M{"rating_ids": ratingID}}
	result, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": albumID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("album not found")
	}

	return nil
}

func (r *albumRepository) Delete(id string) error {
	result, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("album not found")
	}

	return nil
}
