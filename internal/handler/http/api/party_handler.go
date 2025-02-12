package api

import (
	"encoding/json"
	"net/http"

	"vinyl-party/internal/dto"
	party_mapper "vinyl-party/internal/mapper/custom/party"
	"vinyl-party/internal/service"

	"github.com/go-chi/chi/v5"
)

type PartyHandler struct {
	partyService service.PartyService
}

func NewPartyHandler(partyService service.PartyService) *PartyHandler {
	return &PartyHandler{
		partyService: partyService,
	}
}

func (h *PartyHandler) CreateParty(w http.ResponseWriter, r *http.Request) {
	var req dto.PartyCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	party := party_mapper.CreateDTOToEntity(req)
	err := h.partyService.Create(&party)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(party)
}

func (h *PartyHandler) GetPartyInfo(w http.ResponseWriter, r *http.Request) {
	partyID := chi.URLParam(r, "id")
	if partyID == "" {
		http.Error(w, "Отсутствует идентификатор вечеринки", http.StatusBadRequest)
		return
	}

	partyInfo, err := h.partyService.GetByID(partyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(partyInfo)
}

func (h *PartyHandler) AddAlbumToParty(w http.ResponseWriter, r *http.Request) {
	partyID := chi.URLParam(r, "id")
	if partyID == "" {
		http.Error(w, "Отсутствует идентификатор вечеринки", http.StatusBadRequest)
		return
	}

	var input struct {
		AlbumID string `json:"album_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if err := h.partyService.AddAlbum(partyID, input.AlbumID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
