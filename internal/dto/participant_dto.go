package dto

type CreateParticipantDTO struct {
	UserID  string `json:"user_id"`
	PartyID string `json:"party_id"`
}

type ParticipantInfoDTO struct {
	ID      string           `json:"id"`
	User    UserShortInfoDTO `json:"user"`
	PartyID string           `json:"party_id"`
	Role    PartyRoleInfoDTO `json:"role"`
}
