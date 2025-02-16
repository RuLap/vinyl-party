package service

import (
	"context"
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"

	"github.com/google/uuid"
)

type AlbumService interface {
	Create(ctx context.Context, album *entity.Album) (string, error)
	GetByID(ctx context.Context, id string) (*entity.Album, error)
	GetByIDs(ctx context.Context, ids []string) ([]*entity.Album, error)
	AddRating(ctx context.Context, albumID string, ratingID string) error
	Delete(ctx context.Context, id string) error
}

type albumService struct {
	albumRepo repository.AlbumRepository
}

func NewAlbumService(albumRepo repository.AlbumRepository) AlbumService {
	return &albumService{albumRepo: albumRepo}
}

func (s *albumService) Create(ctx context.Context, album *entity.Album) (string, error) {
	album.ID = uuid.NewString()
	err := s.albumRepo.Create(ctx, album)
	if err != nil {
		return "", err
	}

	return album.ID, nil
}

func (s *albumService) GetByID(ctx context.Context, id string) (*entity.Album, error) {
	return s.albumRepo.GetByID(ctx, id)
}

func (s *albumService) GetByIDs(ctx context.Context, ids []string) ([]*entity.Album, error) {
	return s.albumRepo.GetByIDs(ctx, ids)
}

func (s *albumService) AddRating(ctx context.Context, albumID string, ratingID string) error {
	return s.albumRepo.AddRating(ctx, albumID, ratingID)
}

func (s *albumService) Delete(ctx context.Context, id string) error {
	return s.albumRepo.Delete(ctx, id)
}
