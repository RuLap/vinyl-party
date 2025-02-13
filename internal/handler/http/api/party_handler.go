package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"vinyl-party/internal/entity"

	"vinyl-party/internal/dto"
	album_mapper "vinyl-party/internal/mapper/custom/album"
	participant_mapper "vinyl-party/internal/mapper/custom/participant"
	party_mapper "vinyl-party/internal/mapper/custom/party"
	party_role_mapper "vinyl-party/internal/mapper/custom/party_role"
	rating_mapper "vinyl-party/internal/mapper/custom/rating"
	user_mapper "vinyl-party/internal/mapper/custom/user"
	"vinyl-party/internal/service"

	"github.com/go-chi/chi/v5"
)

type PartyHandler struct {
	userService        service.UserService
	albumService       service.AlbumService
	partyService       service.PartyService
	ratingService      service.RatingService
	spotifyService     service.SpotifyService
	partyRoleService   service.PartyRoleService
	participantService service.ParticipantService
}

func NewPartyHandler(
	userService service.UserService,
	albumService service.AlbumService,
	partyService service.PartyService,
	ratingService service.RatingService,
	spotifyService service.SpotifyService,
	partyRoleService service.PartyRoleService,
	participantService service.ParticipantService) *PartyHandler {
	return &PartyHandler{
		userService:        userService,
		albumService:       albumService,
		partyService:       partyService,
		ratingService:      ratingService,
		spotifyService:     spotifyService,
		partyRoleService:   partyRoleService,
		participantService: participantService,
	}
}

func (h *PartyHandler) CreateParty(w http.ResponseWriter, r *http.Request) {
	var req dto.PartyCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByID(req.HostID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	adminRole, err := h.partyRoleService.GetByName("Админ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	party := party_mapper.CreateDTOToEntity(req)
	err = h.partyService.Create(&party)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	participant := &entity.Participant{
		ID:      uuid.NewString(),
		UserID:  user.ID,
		RoleID:  adminRole.ID,
		PartyID: party.ID,
	}
	err = h.participantService.Create(participant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(party)
}

func (h *PartyHandler) GetAllParties(w http.ResponseWriter, e *http.Request) {
	parties, err := h.partyService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(parties)
}

func (h *PartyHandler) GetPartyInfo(w http.ResponseWriter, r *http.Request) {
	partyID := chi.URLParam(r, "id")
	if partyID == "" {
		http.Error(w, "Отсутствует идентификатор вечеринки", http.StatusBadRequest)
		return
	}

	party, err := h.partyService.GetByID(partyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	albumDTOs := h.getPartyAlbumDTOs(w, party)
	participantDTOs := h.getPartyParticipantDTOs(w, party)

	partyDTO := party_mapper.EntityToInfoDTO(*party, albumDTOs, participantDTOs)

	json.NewEncoder(w).Encode(partyDTO)
}

func (h *PartyHandler) getPartyAlbumDTOs(w http.ResponseWriter, party *entity.Party) []dto.AlbumInfoDTO {
	albums, err := h.albumService.GetByPartyID(party.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	albumDTOs := make([]dto.AlbumInfoDTO, 0)
	if len(albums) != 0 {
		for _, album := range albums {
			ratings, err := h.ratingService.GetByAlbumID(album.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			ratingDTOs := make([]dto.RatingInfoDTO, 0)
			averageRating := 0
			for _, rating := range ratings {
				rater, err := h.userService.GetByID(rating.UserID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				raterDTO := user_mapper.EntityToShortInfoDTO(*rater)

				albumDTO := rating_mapper.EntityToInfoDTO(*rating, raterDTO)
				ratingDTOs = append(ratingDTOs, albumDTO)
				averageRating += rating.Score
			}

			albumDTO := album_mapper.EntityToInfoDTO(album, ratingDTOs, averageRating)
			albumDTOs = append(albumDTOs, albumDTO)
		}
	}

	return albumDTOs
}

func (h *PartyHandler) getPartyParticipantDTOs(w http.ResponseWriter, party *entity.Party) []dto.ParticipantInfoDTO {
	participants, err := h.participantService.GetByPartyID(party.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	participantDTOs := make([]dto.ParticipantInfoDTO, 0)
	if len(participants) != 0 {
		for _, participant := range participants {
			user, err := h.userService.GetByID(participant.UserID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			userDTO := user_mapper.EntityToShortInfoDTO(*user)
			role, err := h.partyRoleService.GetByID(participant.RoleID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			roleDTO := party_role_mapper.EntityToPartyRoleDTO(*role)
			participantDTO := participant_mapper.EntityToParticipantInfoDTO(*participant, userDTO, roleDTO)
			participantDTOs = append(participantDTOs, participantDTO)
		}
	}

	return participantDTOs
}

func (h *PartyHandler) AddAlbumToParty(w http.ResponseWriter, r *http.Request) {
	partyID := chi.URLParam(r, "id")
	if partyID == "" {
		http.Error(w, "Отсутствует идентификатор вечеринки", http.StatusBadRequest)
		return
	}

	party, err := h.partyService.GetByID(partyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	var req dto.AlbumAddFromSpotifyDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	spotifyAlbum, err := h.spotifyService.GetAlbumFromLink(req.SpotifyUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	albums, err := h.albumService.GetByPartyID(party.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for _, album := range albums {
		if album.SpotifyUrl == spotifyAlbum.Url {
			http.Error(w, "Album already exists", http.StatusBadRequest)
			return
		}
	}

	albumCreateDTO := album_mapper.SpotifyDTOToEntity(spotifyAlbum)
	albumCreateDTO.PartyID = party.ID
	_, err = h.albumService.Create(&albumCreateDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PartyHandler) AddParticipantToParty(w http.ResponseWriter, r *http.Request) {
	partyID := chi.URLParam(r, "id")
	if partyID == "" {
		http.Error(w, "Отсутствует идентификатор вечеринки", http.StatusBadRequest)
		return
	}

	var input struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if err := h.partyService.AddParticipant(partyID, input.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
