package service

import (
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"

	"github.com/google/uuid"
)

type RatingService interface {
	Create(rating *entity.Rating) error
	FindByID(id string) (*entity.Rating, error)
}

type ratingService struct {
	ratingRepo repository.RatingRepository
}

func NewRatingService(ratingRepo repository.RatingRepository) RatingService {
	return &ratingService{ratingRepo: ratingRepo}
}

func (s *ratingService) Create(rating *entity.Rating) error {
	existingRating, err := s.ratingRepo.FindByID(rating.ID)
	if err == nil && existingRating.ID != "" {
		return err
	}

	rating.ID = uuid.NewString()
	return s.ratingRepo.Create(rating)
}

func (s *ratingService) FindByID(id string) (*entity.Rating, error) {
	rating, err := s.ratingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return rating, nil
}
