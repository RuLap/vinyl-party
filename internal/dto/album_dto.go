package dto

type AlbumCreateDTO struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	CoverUrl string `json:"cover_url"`
}

type AlbumShortDTO struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	CoverUrl string `json:"cover_url"`
}

type AlbumInfoDTO struct {
	ID            string          `json:"id"`
	Title         string          `json:"title"`
	Artist        string          `json:"artist"`
	CoverUrl      string          `json:"cover_url"`
	Ratings       []RatingInfoDTO `json:"ratings"`
	AverageRating int             `json:"average_rating"`
}
