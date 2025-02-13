package repository

import (
	"context"
	"vinyl-party/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartyRoleRepository interface {
	Create(role *entity.PartyRole) error
	GetByID(id string) (*entity.PartyRole, error)
	GetByName(name string) (*entity.PartyRole, error)
}

type partyRoleRepository struct {
	collection *mongo.Collection
}

func NewPartyRoleRepository(db *mongo.Database) PartyRoleRepository {
	return &partyRoleRepository{collection: db.Collection("party_roles")}
}

func (r *partyRoleRepository) Create(role *entity.PartyRole) error {
	_, err := r.collection.InsertOne(context.Background(), role)
	return err
}

func (r *partyRoleRepository) GetByID(id string) (*entity.PartyRole, error) {
	var role entity.PartyRole
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&role)
	if err != nil {
		return nil, err
	}

	return &role, err
}

func (r *partyRoleRepository) GetByName(name string) (*entity.PartyRole, error) {
	var role entity.PartyRole
	err := r.collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&role)
	if err != nil {
		return nil, err
	}

	return &role, err
}
