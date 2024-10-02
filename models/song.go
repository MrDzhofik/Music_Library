package models

import "time"

type Song struct {
	ID          *int       `json:"id"`
	Name        *string    `json:"song"`
	Group       *int       `json:"group"`
	Group_name  *string    `json:"group_name"`
	ReleaseDate *time.Time `json:"releaseDate"`
	Text        *string    `json:"text"`
	Link        *string    `json:"link"`
}

type Group struct {
	Name *string `json:"group"`
}
