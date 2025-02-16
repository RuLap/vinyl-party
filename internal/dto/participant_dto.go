package dto

type CreateParticipantDTO struct {
	UserID string `json:"user_id"`
}

type ParticipantInfoDTO struct {
	User UserShortInfoDTO `json:"user"`
	Role string           `json:"role"`
}
