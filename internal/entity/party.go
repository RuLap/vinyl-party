package entity

type Party struct {
	ID              string   `bson:"_id"`
	HostID          string   `bson:"host_id"`
	Title           string   `bson:"title"`
	Description     string   `bson:"description"`
	Date            string   `bson:"date"`
	AlbumsIDs       []string `bson:"album_ids"`
	ParticipantsIDs []string `bson:"participant_ids"`
}
