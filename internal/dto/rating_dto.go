package dto

type RatingCreateDTO struct {
	User  UserShortInfoDTO `json:"user"`
	Album AlbumInfoDTO     `json:"album"`
	Score int              `json:"score"`
}

type RatingInfoDTO struct {
	ID   string           `json:"id"`
	User UserShortInfoDTO `json:"user"`

	Score int `json:"score"`
}
