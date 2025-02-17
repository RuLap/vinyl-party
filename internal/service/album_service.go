package service

import (
	"context"
	"time"
	"vinyl-party/internal/dto"
	album_mapper "vinyl-party/internal/mapper/custom/album"
	rating_mapper "vinyl-party/internal/mapper/custom/rating"
	user_mapper "vinyl-party/internal/mapper/custom/user"
	"vinyl-party/internal/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

type AlbumService interface {
	GetByID(ctx context.Context, id string) (*dto.AlbumInfoDTO, error)
	GetByIDs(ctx context.Context, ids []string) ([]*dto.AlbumInfoDTO, error)
	AddRating(ctx context.Context, albumID string, ratingDTO *dto.RatingCreateDTO) (*dto.RatingInfoDTO, error)
	Delete(ctx context.Context, id string) error
}

type albumService struct {
	albumRepo  repository.AlbumRepository
	ratingRepo repository.RatingRepository
	userRepo   repository.UserRepository
	client     *mongo.Client
}

func NewAlbumService(
	albumRepo repository.AlbumRepository,
	ratingRepo repository.RatingRepository,
	userRepo repository.UserRepository,
	client *mongo.Client) AlbumService {
	return &albumService{albumRepo: albumRepo, ratingRepo: ratingRepo, userRepo: userRepo, client: client}
}

func (s *albumService) GetByID(ctx context.Context, id string) (*dto.AlbumInfoDTO, error) {
	album, err := s.albumRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	ratingDTOs, err := s.getAlbumRatingsByIDs(ctx, album.RatingIDs)
	if err != nil {
		return nil, err
	}

	albumDTO := album_mapper.EntityToInfoDTO(album, ratingDTOs)

	return &albumDTO, nil
}

func (s *albumService) GetByIDs(ctx context.Context, ids []string) ([]*dto.AlbumInfoDTO, error) {
	albums, err := s.albumRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var albumDTOs []*dto.AlbumInfoDTO
	for _, album := range albums {
		ratingDTOs, err := s.getAlbumRatingsByIDs(ctx, album.RatingIDs)
		if err != nil {
			return nil, err
		}

		albumDTO := album_mapper.EntityToInfoDTO(album, ratingDTOs)
		albumDTOs = append(albumDTOs, &albumDTO)
	}

	return albumDTOs, nil
}

func (s *albumService) getAlbumRatingsByIDs(ctx context.Context, ids []string) ([]dto.RatingInfoDTO, error) {
	ratings, err := s.ratingRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var ratingDTOs []dto.RatingInfoDTO
	for _, rating := range ratings {
		userDTO, err := s.getRatingUserByID(ctx, rating.ID)
		if err != nil {
			return nil, err
		}
		ratingDTO := rating_mapper.EntityToInfoDTO(rating, userDTO)

		ratingDTOs = append(ratingDTOs, ratingDTO)
	}

	return ratingDTOs, nil
}

func (s *albumService) getRatingUserByID(ctx context.Context, id string) (*dto.UserShortInfoDTO, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userDTO := user_mapper.EntityToShortInfoDTO(user)

	return &userDTO, nil
}

func (s *albumService) AddRating(ctx context.Context, albumID string, ratingDTO *dto.RatingCreateDTO) (*dto.RatingInfoDTO, error) {
	session, err := s.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	rating := rating_mapper.CreateDTOToEntity(ratingDTO)
	rating.CreatedAt = time.Now().UTC()

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := s.ratingRepo.Create(sc, &rating); err != nil {
			return err
		}

		if err := s.albumRepo.AddRating(sc, albumID, rating.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, rating.UserID)
	if err != nil {
		return nil, err
	}
	userDTO := user_mapper.EntityToShortInfoDTO(user)

	result := rating_mapper.EntityToInfoDTO(&rating, &userDTO)

	return &result, nil
}

func (s *albumService) Delete(ctx context.Context, id string) error {
	return s.albumRepo.Delete(ctx, id)
}
