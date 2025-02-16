package api

import (
	"encoding/json"
	"net/http"
	"time"
	"vinyl-party/internal/entity"

	"github.com/google/uuid"

	"vinyl-party/internal/dto"
	party_mapper "vinyl-party/internal/mapper/custom/party"
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

	if err := h.partyService.Create(&party); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	partyDTO := party_mapper.EntityToShortInfoDTO(party)

	w.WriteHeader(http.StatusCreated)
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

	parties, err := h.partyService.GetUserParties(userID, status)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(parties)
}
