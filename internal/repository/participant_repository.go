package repository

import (
	"context"
	"vinyl-party/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ParticipantRepository interface {
	Create(participant *entity.Participant) error
	GetByID(id string) (*entity.Participant, error)
}

type participantRepository struct {
	collection *mongo.Collection
}

func NewParticipantRepository(db *mongo.Database) ParticipantRepository {
	return &participantRepository{collection: db.Collection("participants")}
}

func (r *participantRepository) Create(participant *entity.Participant) error {
	_, err := r.collection.InsertOne(context.Background(), participant)
	return err
}

func (r *participantRepository) GetByID(id string) (*entity.Participant, error) {
	var participant entity.Participant
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&participant)
	if err != nil {
		return nil, err
	}

	return &participant, nil
}
