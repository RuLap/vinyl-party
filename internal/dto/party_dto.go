package dto

type PartyCreateDTO struct {
	HostID      string `json:"host_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type PartyShortInfoDTO struct {
	ID          string           `json:"id"`
	Host        UserShortInfoDTO `json:"host"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Date        string           `json:"date"`
}

type PartyInfoDTO struct {
	ID           string             `json:"id"`
	Host         UserShortInfoDTO   `json:"host"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	Date         string             `json:"date"`
	Albums       []AlbumInfoDTO     `json:"albums"`
	Participants []UserShortInfoDTO `json:"participants"`
}
