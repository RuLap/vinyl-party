package entity

import "time"

type Album struct {
	ID            string    `bson:"_id"`
	Title         string    `bson:"title"`
	Artist        string    `bson:"artist"`
	CoverUrl      string    `bson:"cover_url"`
	SpotifyID     string    `bson:"spotify_id"`
	SpotifyUrl    string    `bson:"spotify_url"`
	CreatedAt     time.Time `bson:"created_at"`
	RatingIDs     []string  `bson:"rating_ids"`
	AverageRating *int      `bson:"average_rating"`
}
