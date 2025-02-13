package party_role

import (
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func EntityToPartyRoleDTO(entity entity.PartyRole) dto.PartyRoleInfoDTO {
	return dto.PartyRoleInfoDTO{
		ID:   entity.ID,
		Name: entity.Name,
	}
}
