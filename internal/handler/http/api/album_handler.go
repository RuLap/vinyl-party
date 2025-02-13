package api

import (
	"encoding/json"
	"net/http"
	"vinyl-party/internal/dto"

	album_mapper "vinyl-party/internal/mapper/custom/album"
	"vinyl-party/internal/service"

	"github.com/go-chi/chi/v5"
)

type AlbumHandler struct {
	albumService service.AlbumService
}

func NewAlbumHandler(albumService service.AlbumService) *AlbumHandler {
	return &AlbumHandler{albumService: albumService}
}

func (h *AlbumHandler) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	var req dto.AlbumCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	album := album_mapper.CreateDTOToEntity(req)

	if _, err := h.albumService.Create(&album); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(album)
}

func (h *AlbumHandler) GetAlbumByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing album id", http.StatusBadRequest)
		return
	}

	album, err := h.albumService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	albumDTO := album_mapper.EntityToShortInfoDTO(*album)

	json.NewEncoder(w).Encode(albumDTO)
}

func (h *AlbumHandler) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing album id", http.StatusBadRequest)
		return
	}

	if err := h.albumService.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AlbumHandler) AddRating(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing album id", http.StatusBadRequest)
		return
	}

	var payload struct {
		RatingID string `json:"rating_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if payload.RatingID == "" {
		http.Error(w, "Missing rating id", http.StatusBadRequest)
		return
	}

	if err := h.albumService.AddRating(id, payload.RatingID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
