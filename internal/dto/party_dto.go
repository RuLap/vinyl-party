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
	ID           string             `bson:"id"`
	Host         UserShortInfoDTO   `json:"host"`
	Title        string             `bson:"title"`
	Description  string             `bson:"description"`
	Date         string             `bson:"date"`
	Albums       []AlbumInfoDTO     `bson:"albums"`
	Participants []UserShortInfoDTO `bson:"participants"`
}
