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
	GetByIDs(id []string) ([]*entity.Party, error)
	GetActiveByIDs(ids []string) ([]*entity.Party, error)
	GetArchiveByIDs(ids []string) ([]*entity.Party, error)
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
	return s.partyRepo.Create(party)
}

func (s *partyService) GetAll() ([]*entity.Party, error) {
	return s.partyRepo.GetAll()
}

func (s *partyService) GetByID(id string) (*entity.Party, error) {
	return s.partyRepo.GetByID(id)
}

func (s *partyService) GetByIDs(ids []string) ([]*entity.Party, error) {
	return s.partyRepo.GetByIDs(ids)
}

func (s *partyService) GetActiveByIDs(ids []string) ([]*entity.Party, error) {
	return s.partyRepo.GetActiveByIDs(ids)
}

func (s *partyService) GetArchiveByIDs(ids []string) ([]*entity.Party, error) {
	return s.partyRepo.GetArchiveByIDs(ids)
}

func (s *partyService) AddAlbum(partyID string, albumID string) error {
	return s.partyRepo.AddAlbum(partyID, albumID)
}

func (s *partyService) AddParticipant(partyID string, userID string) error {
	return s.partyRepo.AddParticipant(partyID, userID)
}
