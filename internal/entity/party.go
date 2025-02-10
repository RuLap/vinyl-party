package entity

type Party struct {
	ID           string   `bson:"_id"`
	HostID       string   `bson:"host_id"`
	Title        string   `bson:"title"`
	Description  string   `bson:"description"`
	Date         string   `bson:"date"`
	Albums       []string `bson:"albums"`
	Participants []string `bson:"participants"`
}
