package dto

import "vinyl-party/internal/entity"

type RatingCreateDTO struct {
	User  entity.User `json:"user"`
	Score int         `json:"score"`
}

type RatingInfoDTO struct {
	ID    string           `json:"id"`
	User  UserShortInfoDTO `json:"user"`
	Score int              `json:"score"`
}
