package repository

import (
	"context"
	"fmt"
	"time"
	"vinyl-party/internal/entity"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartyRepository interface {
	EnsureIndexes() error
	Create(ctx context.Context, party *entity.Party) error
	GetPartiesByUserID(ctx context.Context, userID string, status entity.PartyStatus) ([]*entity.Party, error)
	GetByID(ctx context.Context, id string) (*entity.Party, error)
	AddAlbumToParty(ctx context.Context, partyID, albumID string) error
	AddParticipant(ctx context.Context, partyID string, participant *entity.Participant) error
}

type partyRepository struct {
	collection *mongo.Collection
}

func NewPartyRepository(db *mongo.Database) PartyRepository {
	return &partyRepository{
		collection: db.Collection("parties"),
	}
}

func (r *partyRepository) EnsureIndexes() error {
	participantIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "participants.user_id", Value: 1},
			{Key: "date", Value: 1},
		},
		Options: options.Index().SetName("user_parties"),
	}

	uniqueParticipantIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "_id", Value: 1},
			{Key: "participants.user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	indexes := []mongo.IndexModel{
		participantIndex,
		uniqueParticipantIndex,
	}
	_, err := r.collection.Indexes().CreateMany(context.Background(), indexes)

	return err
}

func (r *partyRepository) Create(ctx context.Context, party *entity.Party) error {
	_, err := r.collection.InsertOne(ctx, party)
	return err
}

func (r *partyRepository) GetPartiesByUserID(ctx context.Context, userID string, status entity.PartyStatus) ([]*entity.Party, error) {
	filter := bson.M{
		"participants.user_id": userID,
	}

	now := time.Now().UTC()
	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)

	switch status {
	case entity.PartyStatusActive:
		filter["date"] = bson.M{"$lte": endOfToday}
	case entity.PartyStatusArchive:
		filter["date"] = bson.M{"$gte": endOfToday}
	}

	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find parties: %v", err)
	}

	var parties []*entity.Party
	if err = cursor.All(ctx, &parties); err != nil {
		return nil, fmt.Errorf("failed to decode parties: %v", err)
	}

	return parties, nil
}

func (r *partyRepository) GetByID(ctx context.Context, id string) (*entity.Party, error) {
	var party entity.Party
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&party)
	if err != nil {
		return nil, err
	}

	return &party, nil
}

func (r *partyRepository) AddAlbumToParty(ctx context.Context, partyID, albumID string) error {
	filter := bson.M{"_id": partyID}
	update := bson.M{"$addToSet": bson.M{"album_ids": albumID}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *partyRepository) AddParticipant(ctx context.Context, partyID string, participant *entity.Participant) error {
	filter := bson.M{"_id": partyID}
	update := bson.M{"$addToSet": bson.M{"participants": participant}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
