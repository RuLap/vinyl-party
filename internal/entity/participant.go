package entity

type Participant struct {
	ID      string `bson:"_id"`
	UserID  string `bson:"user_id"`
	PartyID string `bson:"party_id"`
	RoleID  string `bson:"role_id"`
}
