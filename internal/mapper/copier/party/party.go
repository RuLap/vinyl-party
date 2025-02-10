package party

import (
	"github.com/jinzhu/copier"
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func CreateDTOToEntity(dto dto.PartyCreateDTO) (entity.Party, error) {
	var party entity.Party
	err := copier.Copy(&party, &dto)
	if err != nil {
		return entity.Party{}, err
	}

	return party, nil
}

func EntityToShortInfoDTO(party entity.Party) (dto.PartyShortInfoDTO, error) {
	var shortInfoDTO dto.PartyShortInfoDTO
	err := copier.Copy(&shortInfoDTO, &party)
	if err != nil {
		return dto.PartyShortInfoDTO{}, err
	}

	return shortInfoDTO, nil
}

func EntityToInfoDTO(party entity.Party) (dto.PartyInfoDTO, error) {
	var infoDTO dto.PartyInfoDTO
	err := copier.Copy(&infoDTO, &party)
	if err != nil {
		return dto.PartyInfoDTO{}, err
	}

	return infoDTO, nil
}
