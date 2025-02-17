package api

import (
	"encoding/json"
	"net/http"
	"vinyl-party/internal/entity"

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
	var req *dto.PartyCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	createdParty, err := h.partyService.Create(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdParty)
}

func (h *PartyHandler) GetParty(w http.ResponseWriter, r *http.Request) {
	partyID := chi.URLParam(r, "id")
	if partyID == "" {
		http.Error(w, "PartyID не найден", http.StatusBadRequest)
		return
	}

	party, err := h.partyService.GetByID(r.Context(), partyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(party)
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
	partyID := chi.URLParam(r, "id")

	var req *dto.AlbumAddFromSpotifyDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	createAlbumDTO, err := h.spotifyService.GetAlbumFromLink(req.SpotifyUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdAlbumDTO, err := h.partyService.AddAlbum(r.Context(), partyID, createAlbumDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdAlbumDTO)
}

func (h *PartyHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	partyID := chi.URLParam(r, "id")

	var req *dto.ParticipantCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	createdParticipantDTO, err := h.partyService.AddParticipant(r.Context(), partyID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdParticipantDTO)
}
