package entity

import "time"

type Party struct {
	ID          string    `bson:"_id"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Date        time.Time `bson:"date"`
}
