package slide

import "time"

// Info represents an individual Slide.
type Info struct {
	ID          string    `db:"slide_id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Link        string    `db:"link" json:"link"`
	Image       string    `db:"image" json:"image"`
	SubTitle    string    `db:"sub_title" json:"sub_title"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewSlide contains information needed to create a new Slide.
type NewSlide struct {
	Title    string `json:"title"  validate:"required"`
	Link     string `json:"link"`
	SubTitle string `json:"sub_title"`
	Image    string `json:"image"`
}

// UpdateSlide in database
type UpdateSlide struct {
	Title    *string `json:"title"  validate:"required"`
	Link     *string `json:"link"`
	SubTitle *string `json:"sub_title"`
	Image    *string `json:"image"`
}
