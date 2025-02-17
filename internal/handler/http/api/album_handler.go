package api

import (
	"encoding/json"
	"net/http"

	"vinyl-party/internal/dto"
	"vinyl-party/internal/service"

	"github.com/go-chi/chi/v5"
)

type AlbumHandler struct {
	albumService service.AlbumService
}

func NewAlbumHandler(albumService service.AlbumService) *AlbumHandler {
	return &AlbumHandler{albumService: albumService}
}

func (h *AlbumHandler) AddRating(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing album id", http.StatusBadRequest)
		return
	}

	var req *dto.RatingCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	album, err := h.albumService.AddRating(r.Context(), id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(album)
}
