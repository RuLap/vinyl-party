package api

import (
	"encoding/json"
	"net/http"
	"time"
	"vinyl-party/internal/dto"
	user_mapper "vinyl-party/internal/mapper/custom/user"
	"vinyl-party/internal/service"
	"vinyl-party/pkg/jwt_helper"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := user_mapper.RegisterDTOToEntity(req)
	user.CreatedAt = time.Now()

	if err := h.userService.Register(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.UserLoginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := jwt_helper.GenerateJWT(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   3600 * 24 * 7,
	})

	loginResponseDTO := user_mapper.EntityToLoginRepsponseDTO(*user, token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponseDTO)
}
