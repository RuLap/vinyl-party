package service

import (
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"
)

type ParticipantService interface {
	Create(participant *entity.Participant) error
	GetByID(id string) (*entity.Participant, error)
}

type participantService struct {
	participantRepo repository.ParticipantRepository
}

func NewParticipantService(participantRepo repository.ParticipantRepository) ParticipantService {
	return &participantService{participantRepo: participantRepo}
}

func (s *participantService) Create(participant *entity.Participant) error {
	return s.participantRepo.Create(participant)
}

func (s *participantService) GetByID(id string) (*entity.Participant, error) {
	return s.participantRepo.GetByID(id)
}
