package entity

type Rating struct {
	ID      string `bson:"_id"`
	UserID  string `bson:"user_id"`
	AlbumID string `bson:"album_id"`
	Score   int    `bson:"score"`
}
