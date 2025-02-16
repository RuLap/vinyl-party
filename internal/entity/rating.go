package entity

import "time"

type Rating struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	Score     int       `bson:"score"`
	CreatedAt time.Time `bson:"created_at"`
}
