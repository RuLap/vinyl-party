package repository

import (
	"context"
	"errors"
	"time"
	"vinyl-party/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartyRepository interface {
	Create(party *entity.Party) error
	GetAll() ([]*entity.Party, error)
	GetByID(id string) (*entity.Party, error)
	GetByIDs(ids []string) ([]*entity.Party, error)
	GetActiveByIDs(ids []string) ([]*entity.Party, error)
	GetArchiveByIDs(ids []string) ([]*entity.Party, error)
	AddAlbum(partyID string, albumID string) error
	AddParticipant(partyID string, userID string) error
}

type partyRepository struct {
	collection *mongo.Collection
}

func NewPartyRepository(db *mongo.Database) PartyRepository {
	return &partyRepository{
		collection: db.Collection("parties"),
	}
}

func (r *partyRepository) Create(party *entity.Party) error {
	_, err := r.collection.InsertOne(context.Background(), party)
	return err
}

func (r *partyRepository) GetAll() ([]*entity.Party, error) {
	var parties []*entity.Party
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.Background(), &parties); err != nil {
		return nil, err
	}

	return parties, nil
}

func (r *partyRepository) GetByID(id string) (*entity.Party, error) {
	var party entity.Party
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&party)
	if err != nil {
		return nil, err
	}

	return &party, nil
}

func (r *partyRepository) GetByIDs(id []string) ([]*entity.Party, error) {
	var parties []*entity.Party
	filter := bson.M{"_id": bson.M{"$in": id}}
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &parties)
	return parties, err
}

func (r *partyRepository) GetActiveByIDs(id []string) ([]*entity.Party, error) {
	var parties []*entity.Party
	now := time.Now()
	filter := bson.M{
		"_id":  bson.M{"$in": id},
		"date": bson.M{"$gt": now},
	}
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &parties)
	return parties, err
}

func (r *partyRepository) GetArchiveByIDs(id []string) ([]*entity.Party, error) {
	var parties []*entity.Party
	now := time.Now()
	filter := bson.M{
		"_id":  bson.M{"$in": id},
		"date": bson.M{"$lt": now},
	}
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &parties)
	return parties, err
}

func (r *partyRepository) AddAlbum(partyID string, albumID string) error {
	update := bson.M{"$addToSet": bson.M{"album_ids": albumID}}
	result, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": partyID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("party not found")
	}

	return nil
}

func (r *partyRepository) AddParticipant(partyID string, userID string) error {
	update := bson.M{"$addToSet": bson.M{"participant_ids": userID}}
	result, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": partyID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("party not found")
	}

	return nil
}
