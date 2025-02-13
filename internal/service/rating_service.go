package service

import (
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"

	"github.com/google/uuid"
)

type RatingService interface {
	Create(rating *entity.Rating) error
	GetByID(id string) (*entity.Rating, error)
	GetByAlbumID(albumID string) ([]*entity.Rating, error)
}

type ratingService struct {
	ratingRepo repository.RatingRepository
}

func NewRatingService(ratingRepo repository.RatingRepository) RatingService {
	return &ratingService{ratingRepo: ratingRepo}
}

func (s *ratingService) Create(rating *entity.Rating) error {
	existingRating, err := s.ratingRepo.GetByID(rating.ID)
	if err == nil && existingRating.ID != "" {
		return err
	}

	rating.ID = uuid.NewString()
	return s.ratingRepo.Create(rating)
}

func (s *ratingService) GetByID(id string) (*entity.Rating, error) {
	return s.ratingRepo.GetByID(id)
}

func (s *ratingService) GetByAlbumID(albumID string) ([]*entity.Rating, error) {
	return s.ratingRepo.GetByAlbumID(albumID)
}
