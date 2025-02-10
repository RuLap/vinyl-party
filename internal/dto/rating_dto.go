package dto

type RatingCreateDTO struct {
	UserID string `bson:"user_id"`
	Score  int    `bson:"score"`
}

type RatingInfoDTO struct {
	User  UserShortInfoDTO `bson:"user"`
	Score int              `bson:"score"`
}
