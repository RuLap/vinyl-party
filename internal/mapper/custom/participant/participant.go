package participant

import (
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func EntityToParticipantInfoDTO(entity entity.Participant, user dto.UserShortInfoDTO, role dto.PartyRoleInfoDTO) dto.ParticipantInfoDTO {
	return dto.ParticipantInfoDTO{
		ID:      entity.ID,
		User:    user,
		PartyID: entity.PartyID,
		Role:    role,
	}
}
