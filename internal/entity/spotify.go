package entity

type SpotifyAlbum struct {
	Name          string             `json:"name"`
	Artists       []SpotifyArtist    `json:"artists"`
	Images        []SpotifyItemImage `json:"images"`
	ArtistsString string             `json:"artists_string";omitempty`
	CoverUrl      string             `json:"cover_url";omitempty`
	Url           string             `json:"url";omitempty`
}

type SpotifyAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type SpotifyArtist struct {
	Name string `json:"name"`
}

type SpotifyItemImage struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}
