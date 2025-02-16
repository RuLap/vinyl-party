package dto

type RatingCreateDTO struct {
	UserID string `json:"user_id"`
	Score  int    `json:"score"`
}

type RatingInfoDTO struct {
	User  UserShortInfoDTO `json:"user"`
	Score int              `json:"score"`
}
