package api

import (
	"encoding/json"
	"net/http"
	"vinyl-party/internal/entity"

	"vinyl-party/internal/dto"
	album_mapper "vinyl-party/internal/mapper/custom/album"
	party_mapper "vinyl-party/internal/mapper/custom/party"
	rating_mapper "vinyl-party/internal/mapper/custom/rating"
	user_mapper "vinyl-party/internal/mapper/custom/user"
	"vinyl-party/internal/service"

	"github.com/go-chi/chi/v5"
)

type PartyHandler struct {
	partyService   service.PartyService
	userService    service.UserService
	albumService   service.AlbumService
	ratingService  service.RatingService
	spotifyService service.SpotifyService
}

func NewPartyHandler(
	partyService service.PartyService,
	userService service.UserService,
	albumService service.AlbumService,
	ratingService service.RatingService,
	spotifyService service.SpotifyService) *PartyHandler {
	return &PartyHandler{
		partyService:   partyService,
		userService:    userService,
		albumService:   albumService,
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

	party := party_mapper.CreateDTOToEntity(req)
	err := h.partyService.Create(&party)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	host, err := h.userService.GetByID(party.HostID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	hostDTO := user_mapper.EntityToShortInfoDTO(*host)
	albumDTOs := h.getPartyAlbumDTOs(w, party)
	participantDTOs := h.getPartyParticipantDTOs(w, party)

	partyDTO := party_mapper.EntityToInfoDTO(*party, hostDTO, albumDTOs, participantDTOs)

	json.NewEncoder(w).Encode(partyDTO)
}

func (h *PartyHandler) getPartyAlbumDTOs(w http.ResponseWriter, party *entity.Party) []dto.AlbumInfoDTO {
	albums, err := h.albumService.GetByIDs(party.AlbumsIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	albumDTOs := make([]dto.AlbumInfoDTO, 0)
	if len(albums) != 0 {
		for _, album := range albums {
			ratings, err := h.ratingService.GetByIDs(album.RaitngIDs)
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

func (h *PartyHandler) getPartyParticipantDTOs(w http.ResponseWriter, party *entity.Party) []dto.UserShortInfoDTO {
	participants, err := h.userService.GetByIDs(party.ParticipantsIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	participantDTOs := make([]dto.UserShortInfoDTO, 0)
	if len(participants) != 0 {
		for _, participant := range participants {
			participantDTO := user_mapper.EntityToShortInfoDTO(*participant)
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

	albums, err := h.albumService.GetByIDs(party.AlbumsIDs)
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
	insertedID, err := h.albumService.Create(&albumCreateDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := h.partyService.AddAlbum(partyID, insertedID); err != nil {
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
