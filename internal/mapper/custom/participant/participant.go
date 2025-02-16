package participant

import (
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func CreateDTOToEntity(dto dto.CreateParticipantDTO) entity.Participant {
	return entity.Participant{
		UserID: dto.UserID,
	}
}
