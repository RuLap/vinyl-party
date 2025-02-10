package entity

type Rating struct {
	UserID string `bson:"user_id"`
	Score  int    `bson:"score"`
}
