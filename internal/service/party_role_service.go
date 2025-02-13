package service

import (
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"
)

type PartyRoleService interface {
	Create(role *entity.PartyRole) error
	GetByID(id string) (*entity.PartyRole, error)
}

type partyRoleService struct {
	partyRoleRepo repository.PartyRoleRepository
}

func NewPartyRoleService(partyRoleRepo repository.PartyRoleRepository) PartyRoleService {
	return &partyRoleService{partyRoleRepo: partyRoleRepo}
}

func (s *partyRoleService) Create(role *entity.PartyRole) error {
	return s.partyRoleRepo.Create(role)
}

func (s *partyRoleService) GetByID(id string) (*entity.PartyRole, error) {
	return s.partyRoleRepo.GetByID(id)
}
