package service

import (
	"context"
	"time"
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
	album_mapper "vinyl-party/internal/mapper/custom/album"
	participant_mapper "vinyl-party/internal/mapper/custom/participant"
	party_mapper "vinyl-party/internal/mapper/custom/party"
	rating_mapper "vinyl-party/internal/mapper/custom/rating"
	user_mapper "vinyl-party/internal/mapper/custom/user"
	"vinyl-party/internal/repository"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"
)

type PartyService interface {
	Create(ctx context.Context, party *dto.PartyCreateDTO) (*dto.PartyShortInfoDTO, error)
	GetByID(ctx context.Context, id string) (*dto.PartyInfoDTO, error)
	GetUserParties(ctx context.Context, userID string, status entity.PartyStatus) ([]*dto.PartyShortInfoDTO, error)
	AddAlbum(ctx context.Context, partyID string, album *entity.SpotifyAlbum) (*dto.AlbumInfoDTO, error)
	AddParticipant(ctx context.Context, partyID string, participantDTO *dto.ParticipantCreateDTO) (*dto.ParticipantInfoDTO, error)
}

type partyService struct {
	partyRepo  repository.PartyRepository
	albumRepo  repository.AlbumRepository
	ratingRepo repository.RatingRepository
	userRepo   repository.UserRepository
	client     *mongo.Client
}

func NewPartyService(
	partyRepo repository.PartyRepository,
	albumRepo repository.AlbumRepository,
	ratingRepo repository.RatingRepository,
	userRepo repository.UserRepository,
	client *mongo.Client) PartyService {
	return &partyService{partyRepo: partyRepo, albumRepo: albumRepo, ratingRepo: ratingRepo, userRepo: userRepo, client: client}
}

func (s *partyService) Create(ctx context.Context, partyDTO *dto.PartyCreateDTO) (*dto.PartyShortInfoDTO, error) {
	participant := entity.Participant{
		UserID: partyDTO.HostID,
		Role:   entity.PartyRoleAdmin,
	}

	party := party_mapper.CreateDTOToEntity(partyDTO)
	party.ID = uuid.NewString()
	party.CreatedAt = time.Now()
	party.AlbumIDs = make([]string, 0)
	party.Participants = append(party.Participants, participant)

	err := s.partyRepo.Create(ctx, &party)
	if err != nil {
		return nil, err
	}

	result := party_mapper.EntityToShortInfoDTO(&party)
	return &result, nil
}

func (s *partyService) GetByID(ctx context.Context, id string) (*dto.PartyInfoDTO, error) {
	party, err := s.partyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	albumDTOs, err := s.getAlbumsByIDs(ctx, party.AlbumIDs)
	if err != nil {
		return nil, err
	}

	participantDTOs, err := s.getParticipantDTOs(ctx, party.Participants)
	if err != nil {
		return nil, err
	}

	partyDTO := party_mapper.EntityToInfoDTO(party, albumDTOs, participantDTOs)

	return &partyDTO, nil
}

func (s *partyService) getAlbumsByIDs(ctx context.Context, ids []string) ([]dto.AlbumInfoDTO, error) {
	albums, err := s.albumRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var albumDTOs []dto.AlbumInfoDTO
	for _, album := range albums {
		ratingDTOs, err := s.getAlbumRatingsByIDs(ctx, album.RatingIDs)
		if err != nil {
			return nil, err
		}

		albumDto := album_mapper.EntityToInfoDTO(album, ratingDTOs)
		albumDTOs = append(albumDTOs, albumDto)
	}

	return albumDTOs, nil
}

func (s *partyService) getAlbumRatingsByIDs(ctx context.Context, ids []string) ([]dto.RatingInfoDTO, error) {
	ratings, err := s.ratingRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var ratingDTOs []dto.RatingInfoDTO
	for _, rating := range ratings {
		user, err := s.userRepo.GetByID(ctx, rating.UserID)
		if err != nil {
			return nil, err
		}

		userDTO := user_mapper.EntityToShortInfoDTO(user)
		ratingDTO := rating_mapper.EntityToInfoDTO(rating, &userDTO)

		ratingDTOs = append(ratingDTOs, ratingDTO)
	}

	return ratingDTOs, nil
}

func (s *partyService) getParticipantDTOs(ctx context.Context, participants []entity.Participant) ([]dto.ParticipantInfoDTO, error) {
	var participantDTOs []dto.ParticipantInfoDTO
	for _, participant := range participants {
		user, err := s.userRepo.GetByID(ctx, participant.UserID)
		if err != nil {
			return nil, err
		}

		userDTO := user_mapper.EntityToShortInfoDTO(user)

		participantDTO := participant_mapper.EntityToParticipantInfoDTO(participant, userDTO)
		participantDTOs = append(participantDTOs, participantDTO)
	}

	return participantDTOs, nil
}

func (s *partyService) GetUserParties(ctx context.Context, userID string, status entity.PartyStatus) ([]*dto.PartyShortInfoDTO, error) {

	parties, err := s.partyRepo.GetPartiesByUserID(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	var partyDTOs []*dto.PartyShortInfoDTO
	for _, party := range parties {
		partyDTO := party_mapper.EntityToShortInfoDTO(party)
		partyDTOs = append(partyDTOs, &partyDTO)
	}

	return partyDTOs, nil
}

func (s *partyService) AddAlbum(ctx context.Context, partyID string, spotifyAlbum *entity.SpotifyAlbum) (*dto.AlbumInfoDTO, error) {
	session, err := s.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	album := album_mapper.SpotifyDTOToEntity(spotifyAlbum)
	album.ID = uuid.NewString()
	album.RatingIDs = make([]string, 0)
	album.CreatedAt = time.Now().UTC()

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := s.albumRepo.Create(sc, &album); err != nil {
			return err
		}

		if err := s.partyRepo.AddAlbumToParty(sc, partyID, album.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	ratingDTOs := make([]dto.RatingInfoDTO, 0)
	result := album_mapper.EntityToInfoDTO(&album, ratingDTOs)

	return &result, nil
}

func (s *partyService) AddParticipant(ctx context.Context, partyID string, participantDTO *dto.ParticipantCreateDTO) (*dto.ParticipantInfoDTO, error) {
	session, err := s.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	participant := entity.Participant{
		UserID: participantDTO.UserID,
		Role:   entity.PartyRoleGuest,
	}
	participant.CreatedAt = time.Now().UTC()

	if err != nil {
		return nil, err
	}

	err = s.partyRepo.AddParticipant(ctx, partyID, &participant)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, participant.UserID)
	if err != nil {
		return nil, err
	}
	userDTO := user_mapper.EntityToShortInfoDTO(user)

	result := participant_mapper.EntityToParticipantInfoDTO(participant, userDTO)

	return &result, nil
}
