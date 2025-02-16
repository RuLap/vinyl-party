package service

import (
	"context"
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"

	"github.com/google/uuid"
)

type RatingService interface {
	Create(ctx context.Context, rating *entity.Rating) error
	GetByID(ctx context.Context, id string) (*entity.Rating, error)
	GetByIDs(ctx context.Context, ids []string) ([]*entity.Rating, error)
	GetByAlbumID(ctx context.Context, albumID string) ([]*entity.Rating, error)
}

type ratingService struct {
	ratingRepo repository.RatingRepository
}

func NewRatingService(ratingRepo repository.RatingRepository) RatingService {
	return &ratingService{ratingRepo: ratingRepo}
}

func (s *ratingService) Create(ctx context.Context, rating *entity.Rating) error {
	existingRating, err := s.ratingRepo.GetByID(ctx, rating.ID)
	if err == nil && existingRating.ID != "" {
		return err
	}

	rating.ID = uuid.NewString()
	return s.ratingRepo.Create(ctx, rating)
}

func (s *ratingService) GetByID(ctx context.Context, id string) (*entity.Rating, error) {
	return s.ratingRepo.GetByID(ctx, id)
}

func (s *ratingService) GetByAlbumID(ctx context.Context, albumID string) ([]*entity.Rating, error) {
	return s.ratingRepo.GetByAlbumID(ctx, albumID)
}

func (s *ratingService) GetByIDs(ctx context.Context, ids []string) ([]*entity.Rating, error) {
	return s.ratingRepo.GetByIDs(ctx, ids)
}
