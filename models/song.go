package models

type Song struct {
	Name        *string `json:"song"`
	Group       *int    `json:"group"`
	Group_name  *string `json:"group_name"`
	ReleaseDate *string `json:"releaseDate"`
	Text        *string `json:"text"`
	Link        *string `json:"link"`
}

type Group struct {
	Name *string `json:"group"`
}
