package service

import (
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"

	"github.com/google/uuid"
)

type AlbumService interface {
	Create(album *entity.Album) (string, error)
	GetByID(id string) (*entity.Album, error)
	AddRating(albumID string, ratingID string) error
	Delete(id string) error
}

type albumService struct {
	albumRepo repository.AlbumRepository
}

func NewAlbumService(albumRepo repository.AlbumRepository) AlbumService {
	return &albumService{albumRepo: albumRepo}
}

func (s *albumService) Create(album *entity.Album) (string, error) {
	album.ID = uuid.NewString()
	err := s.albumRepo.Create(album)
	if err != nil {
		return "", err
	}

	return album.ID, nil
}

func (s *albumService) GetByID(id string) (*entity.Album, error) {
	return s.albumRepo.GetByID(id)
}

func (s *albumService) AddRating(albumID string, ratingID string) error {
	return s.albumRepo.AddRating(albumID, ratingID)
}

func (s *albumService) Delete(id string) error {
	return s.albumRepo.Delete(id)
}
