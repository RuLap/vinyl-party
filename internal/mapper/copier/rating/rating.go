package rating

import (
	"github.com/jinzhu/copier"
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func CreateDTOToEntity(dto dto.RatingCreateDTO) (entity.Rating, error) {
	var rating entity.Rating
	err := copier.Copy(&rating, &dto)
	if err != nil {
		return entity.Rating{}, err
	}

	return rating, nil
}

func EntityToInfoDTO(rating entity.Rating) (dto.RatingInfoDTO, error) {
	var infoDTO dto.RatingInfoDTO
	err := copier.Copy(&infoDTO, &rating)
	if err != nil {
		return dto.RatingInfoDTO{}, err
	}

	return infoDTO, nil
}
