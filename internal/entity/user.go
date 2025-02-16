package entity

import "time"

type User struct {
	ID        string    `bson:"_id"`
	FirstName string    `bson:"first_name"`
	LastName  string    `bson:"last_name"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	AvatarUrl string    `bson:"avatar_url"`
	CreatedAt time.Time `bson:"created_at"`
}
