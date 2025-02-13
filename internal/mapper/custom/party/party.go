package party

import (
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func CreateDTOToEntity(dto dto.PartyCreateDTO) entity.Party {
	return entity.Party{
		Title:       dto.Title,
		Description: dto.Description,
		Date:        dto.Date,
	}
}

func EntityToShortInfoDTO(party entity.Party) dto.PartyShortInfoDTO {
	return dto.PartyShortInfoDTO{
		ID:          party.ID,
		Title:       party.Title,
		Description: party.Description,
		Date:        party.Date,
	}
}

func EntityToInfoDTO(party entity.Party, albumDTOs []dto.AlbumInfoDTO, participantDTOs []dto.UserShortInfoDTO) dto.PartyInfoDTO {
	return dto.PartyInfoDTO{
		ID:           party.ID,
		Title:        party.Title,
		Description:  party.Description,
		Date:         party.Date,
		Albums:       albumDTOs,
		Participants: participantDTOs,
	}
}
