package dto

type ParticipantInfoDTO struct {
	User UserShortInfoDTO `json:"user"`
	Role string           `json:"role"`
}
