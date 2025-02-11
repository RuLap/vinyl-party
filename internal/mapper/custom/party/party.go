package party

import (
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func CreateDTOToEntity(dto dto.PartyCreateDTO, hostID string) entity.Party {
	return entity.Party{
		HostID:      hostID,
		Title:       dto.Title,
		Description: dto.Description,
		Date:        dto.Date,
	}
}

func EntityToShortInfoDTO(party entity.Party, hostDto dto.UserShortInfoDTO) dto.PartyShortInfoDTO {
	return dto.PartyShortInfoDTO{
		ID:          party.ID,
		Host:        hostDto,
		Title:       party.Title,
		Description: party.Description,
		Date:        party.Date,
	}
}

func EntityToInfoDTO(party entity.Party, host dto.UserShortInfoDTO, albumDtos []dto.AlbumInfoDTO, participantDtos []dto.UserShortInfoDTO) dto.PartyInfoDTO {
	return dto.PartyInfoDTO{
		ID:           party.ID,
		Host:         host,
		Title:        party.Title,
		Description:  party.Description,
		Date:         party.Date,
		Albums:       albumDtos,
		Participants: participantDtos,
	}
}
