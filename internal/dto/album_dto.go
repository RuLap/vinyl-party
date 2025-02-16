package dto

type AlbumCreateDTO struct {
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	CoverUrl   string `json:"cover_url"`
	SpotifyUrl string `json:"spotify_string"`
}

type AlbumInfoDTO struct {
	ID            string          `json:"id"`
	Title         string          `json:"title"`
	Artist        string          `json:"artist"`
	CoverUrl      string          `json:"cover_url"`
	SpotifyUrl    string          `json:"spotify_url"`
	Ratings       []RatingInfoDTO `json:"ratings"`
	AverageRating *int            `json:"average_rating,omitempty"`
}

type AlbumAddFromSpotifyDTO struct {
	SpotifyUrl string `json:"spotify_url"`
}
