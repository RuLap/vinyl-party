package entity

type Album struct {
	ID       string `bson:"_id"`
	Title    string `bson:"title"`
	Artist   string `bson:"artist"`
	CoverUrl string `bson:"cover_url"`
}
