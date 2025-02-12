package entity

type Rating struct {
	ID     string `bson:"_id"`
	UserID string `bson:"user_id"`
	Score  int    `bson:"score"`
}
