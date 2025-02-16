package participant

import (
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func EntityToParticipantInfoDTO(participant entity.Participant, user dto.UserShortInfoDTO) dto.ParticipantInfoDTO {
	return dto.ParticipantInfoDTO{
		User: user,
		Role: string(participant.Role),
	}
}
