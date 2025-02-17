package entity

import "time"

type Participant struct {
	UserID    string    `bson:"user_id"`
	Role      PartyRole `bson:"role"`
	CreatedAt time.Time `bson:"created_at"`
}
