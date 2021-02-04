package categories

import (
	"time"
)

// Info represents an individual Category.
type Info struct {
	ID          string    `db:"category_id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Slug        string    `db:"slug" json:"slug"`
	ParrentID   string    `db:"parrent_id" json:"parrent_id"`
	Description string    `db:"description" json:"description"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewCategory contains information needed to create a new Category.
type NewCategory struct {
	Title       string `json:"title"  validate:"required"`
	Slug        string `json:"slug"  validate:"required"`
	ParrentID   string `json:"parrent_id" `
	Description string `json:"description"`
}

// UpdateCategory in database
type UpdateCategory struct {
	Title       *string `json:"title"  validate:"required"`
	Slug        *string `json:"slug"  validate:"required"`
	ParrentID   *string `json:"parrent_id"`
	Description *string `json:"description"`
}
