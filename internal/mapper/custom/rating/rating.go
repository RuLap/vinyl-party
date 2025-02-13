package rating

import (
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func CreateDTOToEntity(dto dto.RatingCreateDTO) entity.Rating {
	return entity.Rating{
		UserID: dto.UserID,
		Score:  dto.Score,
	}
}

func EntityToInfoDTO(rating entity.Rating, userDTO dto.UserShortInfoDTO) dto.RatingInfoDTO {
	return dto.RatingInfoDTO{
		User:  userDTO,
		Score: rating.Score,
	}
}
