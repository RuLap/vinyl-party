package repository

import "vinyl-party/internal/entity"

type RatingRepository interface {
	Create(rating *entity.Rating) error
	FindByID(id string) (entity.Rating, error)
	FindByUserID(userId string) ([]entity.Rating, error)
}
