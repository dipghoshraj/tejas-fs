package entities

import "time"

type Segment struct {
	ID             string    `json:"id"`
	OrbId          string    `json:"orb_id"`
	SequenceNumber int64     `json:"sequence_number"`
	Checksum       int64     `json:"checksum"`
	Location       string    `json:"location"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
