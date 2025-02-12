package service

import (
	"vinyl-party/internal/entity"
	"vinyl-party/internal/repository"

	"github.com/google/uuid"
)

type AlbumService interface {
	Create(album entity.Album) error
	GetByID(id string) (*entity.Album, error)
}

type albumService struct {
	albumRepo repository.AlbumRepository
}

func NewAlbumService(albumRepo repository.AlbumRepository) AlbumService {
	return &albumService{albumRepo: albumRepo}
}

func (s *albumService) Create(album entity.Album) error {
	album.ID = uuid.NewString()
	return s.albumRepo.Create(&album)
}

func (s *albumService) GetByID(id string) (*entity.Album, error) {
	album, err := s.albumRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return album, nil
}
