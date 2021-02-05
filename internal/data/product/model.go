package product

import (
	"time"
)

// Info represents an individual Product.
type Info struct {
	ID          string    `db:"product_id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Slug        string    `db:"slug" json:"slug"`
	CategoryID  string    `db:"category_id" json:"category_id"`
	Price       float64   `db:"price" json:"price"`
	Description string    `db:"description" json:"description"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// NewProduct contains information needed to create a new Product.
type NewProduct struct {
	Title       string  `json:"title"  validate:"required"`
	Slug        string  `json:"slug"  validate:"required"`
	CategoryID  string  `json:"category_id" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Description string  `json:"description"`
}

// UpdateProduct in database
type UpdateProduct struct {
	Title       *string  `json:"title"  validate:"required"`
	Slug        *string  `json:"slug"  validate:"required"`
	CategoryID  *string  `json:"category_id" validate:"required"`
	Price       *float64 `json:"price" validate:"required"`
	Description *string  `json:"description"`
}
