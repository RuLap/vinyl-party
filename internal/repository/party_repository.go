package repository

import (
	"context"
	"vinyl-party/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartyRepository interface {
	Create(party entity.Party) error
	FindByID(id string) (*entity.Party, error)
}

type partyRepository struct {
	collection *mongo.Collection
}

func NewPartyRepository(db *mongo.Database) PartyRepository {
	return &partyRepository{
		collection: db.Collection("parties"),
	}
}

func (r *partyRepository) Create(party entity.Party) error {
	_, err := r.collection.InsertOne(context.Background(), party)
	return err
}

func (r *partyRepository) FindByID(id string) (*entity.Party, error) {
	var party entity.Party
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&party)
	if err != nil {
		return nil, err
	}

	return &party, nil
}
