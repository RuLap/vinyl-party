package entity

type Party struct {
	ID              string   `bson:"_id"`
	Title           string   `bson:"title"`
	Description     string   `bson:"description"`
	Date            string   `bson:"date"`
	ParticipantsIDs []string `bson:"participant_ids"` // TODO: Move this to participant
}
