package brand

import "time"

// Info represents an individual Brand.
type Info struct {
	ID              string    `db:"brand_id" json:"id"`
	Title           string    `db:"title" json:"title"`
	Slug            string    `db:"slug" json:"slug"`
	Image           string    `db:"image" json:"image"`
	Description     string    `db:"description" json:"description"`
	MetaTitle       string    `db:"meta_title" json:"meta_title"`
	MetaKeywords    string    `db:"meta_keywords" json:"meta_keywords"`
	MetaDescription string    `db:"meta_description" json:"meta_description"`
	DateCreated     time.Time `db:"date_created" json:"date_created"`
	DateUpdated     time.Time `db:"date_updated" json:"date_updated"`
}

// NewBrand contains information needed to create a new Brand.
type NewBrand struct {
	Title           string `json:"title"  validate:"required"`
	Slug            string `json:"slug"  validate:"required"`
	Description     string `json:"description"`
	Image           string `json:"image"`
	MetaTitle       string `json:"meta_title"`
	MetaKeywords    string `json:"meta_keywords"`
	MetaDescription string `json:"meta_description"`
}

// UpdateBrand in database
type UpdateBrand struct {
	Title           *string `json:"title"  validate:"required"`
	Slug            *string `json:"slug"  validate:"required"`
	Description     *string `json:"description"`
	Image           *string `json:"image"`
	MetaTitle       *string `json:"meta_title"`
	MetaKeywords    *string `json:"meta_keywords"`
	MetaDescription *string `json:"meta_description"`
}
