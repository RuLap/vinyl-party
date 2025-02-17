package dto

type ParticipantInfoDTO struct {
	User UserShortInfoDTO `json:"user"`
	Role string           `json:"role"`
}

type ParticipantCreateDTO struct {
	UserID string `json:"user_id"`
}
