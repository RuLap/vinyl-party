package api

import (
	"encoding/json"
	"net/http"

	"vinyl-party/internal/entity"
	"vinyl-party/internal/service"

	"github.com/go-chi/chi/v5"
)

type RatingHandler struct {
	ratingService service.RatingService
}

func NewRatingHandler(ratingService service.RatingService) *RatingHandler {
	return &RatingHandler{
		ratingService: ratingService,
	}
}

func (h *RatingHandler) CreateRating(w http.ResponseWriter, r *http.Request) {
	var rating entity.Rating
	if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if err := h.ratingService.Create(&rating); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rating)
}

func (h *RatingHandler) GetRatingByID(w http.ResponseWriter, r *http.Request) {
	ratingID := chi.URLParam(r, "id")
	if ratingID == "" {
		http.Error(w, "Отсутствует идентификатор рейтинга", http.StatusBadRequest)
		return
	}

	rating, err := h.ratingService.GetByID(ratingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(rating)
}
