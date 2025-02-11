package entity

type Album struct {
	ID         string   `bson:"_id"`
	Title      string   `bson:"title"`
	Artist     string   `bson:"artist"`
	CoverUrl   string   `bson:"cover_url"`
	SpotifyUrl string   `bson:"spotify_url"`
	RaitngIDs  []string `bson:"rating_ids"`
}
