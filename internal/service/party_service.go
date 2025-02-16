package service

import (
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"

	"github.com/google/uuid"
)

type PartyService interface {
	Create(party *entity.Party) error
	GetUserParties(userID string, status entity.PartyStatus) ([]*entity.Party, error)
}

type partyService struct {
	partyRepo repository.PartyRepository
}

func NewPartyService(partyRepo repository.PartyRepository) PartyService {
	return &partyService{partyRepo: partyRepo}
}

func (s *partyService) Create(party *entity.Party) error {
	party.ID = uuid.NewString()
	return s.partyRepo.Create(party)
}

func (s *partyService) GetUserParties(userID string, status entity.PartyStatus) ([]*entity.Party, error) {
	return s.partyRepo.GetPartiesByUserID(userID, status)
}
