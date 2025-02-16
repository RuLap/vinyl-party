package entity

import "time"

type PartyStatus string

type Party struct {
	ID           string        `bson:"_id"`
	Title        string        `bson:"title"`
	Description  string        `bson:"description"`
	Date         time.Time     `bson:"date"`
	CreatedAt    time.Time     `bson:"created_at"`
	AlbumIDs     []string      `bson:"album_ids"`
	Participants []Participant `bson:"participants"`
}

const (
	PartyStatusActive  PartyStatus = "Active"
	PartyStatusArchive PartyStatus = "Archive"
)

func IsValidPartyStatus(s string) bool {
	switch PartyStatus(s) {
	case PartyStatusActive, PartyStatusArchive:
		return true
	default:
		return false
	}
}
