package entity

type Participant struct {
	UserID string    `bson:"user_id"`
	Role   PartyRole `bson:"role"`
}
