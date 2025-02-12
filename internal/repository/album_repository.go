package repository

import (
	"context"
	"vinyl-party/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlbumRepository interface {
	Create(album *entity.Album) error
	GetByID(id string) (*entity.Album, error)
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
