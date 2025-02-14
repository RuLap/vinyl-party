package party

import (
	"time"
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func CreateDTOToEntity(dto dto.PartyCreateDTO) entity.Party {
	date, _ := time.Parse("", dto.Date)
	return entity.Party{
		Title:       dto.Title,
		Description: dto.Description,
		Date:        date,
	}
}

func EntityToShortInfoDTO(party entity.Party) dto.PartyShortInfoDTO {
	return dto.PartyShortInfoDTO{
		ID:          party.ID,
		Title:       party.Title,
		Description: party.Description,
		Date:        party.Date.String(),
	}
}

func EntityToInfoDTO(party entity.Party, albumDTOs []dto.AlbumInfoDTO, participantDTOs []dto.ParticipantInfoDTO) dto.PartyInfoDTO {
	return dto.PartyInfoDTO{
		ID:           party.ID,
		Title:        party.Title,
		Description:  party.Description,
		Date:         party.Date.String(),
		Albums:       albumDTOs,
		Participants: participantDTOs,
	}
}
