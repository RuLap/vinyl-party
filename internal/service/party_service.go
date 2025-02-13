package service

import (
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"

	"github.com/google/uuid"
)

type PartyService interface {
	Create(party *entity.Party) error
	GetAll() ([]*entity.Party, error)
	GetByID(id string) (*entity.Party, error)
	AddAlbum(partyID string, albumID string) error
	AddParticipant(partyID string, userID string) error
}

type partyService struct {
	partyRepo repository.PartyRepository
}

func NewPartyService(partyRepo repository.PartyRepository) PartyService {
	return &partyService{partyRepo: partyRepo}
}

func (s *partyService) Create(party *entity.Party) error {
	party.ID = uuid.NewString()
	party.AlbumsIDs = []string{}
	party.ParticipantsIDs = []string{}
	return s.partyRepo.Create(party)
}

func (s *partyService) GetAll() ([]*entity.Party, error) {
	parties, err := s.partyRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return parties, nil
}

func (s *partyService) GetByID(id string) (*entity.Party, error) {
	party, err := s.partyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return party, nil
}

func (s *partyService) AddAlbum(partyID string, albumID string) error {
	return s.partyRepo.AddAlbum(partyID, albumID)
}

func (s *partyService) AddParticipant(partyID string, userID string) error {
	return s.partyRepo.AddParticipant(partyID, userID)
}
