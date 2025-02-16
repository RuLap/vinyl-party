package service

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"

	"github.com/google/uuid"
)

type PartyService interface {
	Create(ctx context.Context, party *entity.Party) error
	GetByID(ctx context.Context, id string) (*entity.Party, error)
	GetUserParties(ctx context.Context, userID string, status entity.PartyStatus) ([]*entity.Party, error)
	AddAlbum(ctx context.Context, partyID string, album *entity.Album) error
}

type partyService struct {
	partyRepo repository.PartyRepository
	albumRepo repository.AlbumRepository
	client    *mongo.Client
}

func NewPartyService(partyRepo repository.PartyRepository, albumRepo repository.AlbumRepository, client *mongo.Client) PartyService {
	return &partyService{partyRepo: partyRepo}
}

func (s *partyService) Create(ctx context.Context, party *entity.Party) error {
	party.ID = uuid.NewString()
	return s.partyRepo.Create(ctx, party)
}

func (s *partyService) GetByID(ctx context.Context, id string) (*entity.Party, error) {
	return s.partyRepo.GetByID(ctx, id)
}

func (s *partyService) GetUserParties(ctx context.Context, userID string, status entity.PartyStatus) ([]*entity.Party, error) {
	return s.partyRepo.GetPartiesByUserID(ctx, userID, status)
}

func (s *partyService) AddAlbum(ctx context.Context, partyID string, album *entity.Album) error {
	session, err := s.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := s.albumRepo.Create(sc, album); err != nil {
			return err
		}

		if err := s.partyRepo.AddAlbumToParty(sc, partyID, album.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
