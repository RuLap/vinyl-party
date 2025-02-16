package api

import (
	"encoding/json"
	"net/http"
	"time"
	"vinyl-party/internal/entity"
	album_mapper "vinyl-party/internal/mapper/custom/album"
	participant_mapper "vinyl-party/internal/mapper/custom/participant"
	party_mapper "vinyl-party/internal/mapper/custom/party"
	rating_mapper "vinyl-party/internal/mapper/custom/rating"
	user_mapper "vinyl-party/internal/mapper/custom/user"

	"github.com/google/uuid"

	"vinyl-party/internal/dto"
	"vinyl-party/internal/service"

	"github.com/go-chi/chi/v5"
)

type PartyHandler struct {
	userService    service.UserService
	albumService   service.AlbumService
	partyService   service.PartyService
	ratingService  service.RatingService
	spotifyService service.SpotifyService
}

func NewPartyHandler(
	userService service.UserService,
	albumService service.AlbumService,
	partyService service.PartyService,
	ratingService service.RatingService,
	spotifyService service.SpotifyService) *PartyHandler {
	return &PartyHandler{
		userService:    userService,
		albumService:   albumService,
		partyService:   partyService,
		ratingService:  ratingService,
		spotifyService: spotifyService,
	}
}

func (h *PartyHandler) CreateParty(w http.ResponseWriter, r *http.Request) {
	var req dto.PartyCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	participant := entity.Participant{
		UserID: req.HostID,
		Role:   entity.PartyRoleAdmin,
	}

	party := party_mapper.CreateDTOToEntity(req)
	party.ID = uuid.NewString()
	party.CreatedAt = time.Now()
	party.AlbumIDs = make([]string, 0)
	party.Participants = append(party.Participants, participant)

	if err := h.partyService.Create(r.Context(), &party); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	partyDTO := party_mapper.EntityToShortInfoDTO(party)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(partyDTO)
}

func (h *PartyHandler) GetParty(w http.ResponseWriter, r *http.Request) {
	partyID := chi.URLParam(r, "id")
	if partyID == "" {
		http.Error(w, "PartyID not found", http.StatusBadRequest)
		return
	}

	party, err := h.partyService.GetByID(r.Context(), partyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	albums, err := h.albumService.GetByIDs(r.Context(), party.AlbumIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var albumDTOs []dto.AlbumInfoDTO
	for _, album := range albums {
		ratings, err := h.ratingService.GetByIDs(r.Context(), album.RatingIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var ratingDTOs []dto.RatingInfoDTO
		for _, rating := range ratings {
			user, err := h.userService.GetByID(r.Context(), rating.UserID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			userDTO := user_mapper.EntityToShortInfoDTO(*user)
			ratingDTO := rating_mapper.EntityToInfoDTO(*rating, userDTO)

			ratingDTOs = append(ratingDTOs, ratingDTO)
		}

		albumDto := album_mapper.EntityToInfoDTO(album, ratingDTOs)
		albumDTOs = append(albumDTOs, albumDto)
	}
	var participantDTOs []dto.ParticipantInfoDTO
	for _, participant := range party.Participants {
		user, err := h.userService.GetByID(r.Context(), participant.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userDTO := user_mapper.EntityToShortInfoDTO(*user)

		participantDTO := participant_mapper.EntityToParticipantInfoDTO(participant, userDTO)
		participantDTOs = append(participantDTOs, participantDTO)
	}

	partyDTO := party_mapper.EntityToInfoDTO(*party, albumDTOs, participantDTOs)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(partyDTO)
}

func (h *PartyHandler) GetUserParties(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		http.Error(w, "UserID not found", http.StatusBadRequest)
		return
	}

	statusStr := r.URL.Query().Get("status")
	if statusStr != "" && !entity.IsValidPartyStatus(statusStr) {
		http.Error(w, "Invalid party status", http.StatusBadRequest)
		return
	}
	var status entity.PartyStatus
	if statusStr != "" {
		status = entity.PartyStatus(statusStr)
	}

	parties, err := h.partyService.GetUserParties(r.Context(), userID, status)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(parties)
}

func (h *PartyHandler) AddAlbum(w http.ResponseWriter, r *http.Request) {

}
