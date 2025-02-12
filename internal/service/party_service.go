package service

import (
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"

	"github.com/google/uuid"
)

type PartyService interface {
	Create(party *entity.Party) error
	FindByID(id string) (*entity.Party, error)
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

func (s *partyService) FindByID(id string) (*entity.Party, error) {
	party, err := s.partyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return party, nil
}
