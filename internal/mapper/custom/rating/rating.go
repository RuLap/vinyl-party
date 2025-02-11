package rating

import (
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func CreateDTOToEntity(dto dto.RatingCreateDTO, userID string) entity.Rating {
	return entity.Rating{
		UserID: userID,
		Score:  dto.Score,
	}
}

func EntityToInfoDTO(rating entity.Rating, userDto dto.UserShortInfoDTO) dto.RatingInfoDTO {
	return dto.RatingInfoDTO{
		User:  userDto,
		Score: rating.Score,
	}
}
