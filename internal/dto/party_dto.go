package dto

type PartyCreateDTO struct {
	HostID      string `json:"host_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type PartyShortInfoDTO struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type PartyInfoDTO struct {
	ID           string               `json:"id"`
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Date         string               `json:"date"`
	Albums       []AlbumInfoDTO       `json:"albums"`
	Participants []ParticipantInfoDTO `json:"participants"`
}
