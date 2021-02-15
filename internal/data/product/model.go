package product

import (
	"time"
)

// Info represents an individual Product.
type Info struct {
	ID               string    `db:"product_id" json:"id"`
	Title            string    `db:"title" json:"title"`
	Slug             string    `db:"slug" json:"slug"`
	CategoryID       string    `db:"category_id" json:"category_id"`
	BrandID          string    `db:"brand_id" json:"brand_id"`
	Price            float64   `db:"price" json:"price"`
	OldPrice         float64   `db:"old_price" json:"old_price"`
	Image            string    `db:"image" json:"image"`
	ShortDescription string    `db:"short_description" json:"short_description"`
	Description      string    `db:"description" json:"description"`
	MetaTitle        string    `db:"meta_title" json:"meta_title"`
	MetaKeywords     string    `db:"meta_keywords" json:"meta_keywords"`
	MetaDescription  string    `db:"meta_description" json:"meta_description"`
	DateCreated      time.Time `db:"date_created" json:"date_created"`
	DateUpdated      time.Time `db:"date_updated" json:"date_updated"`
}

// NewProduct contains information needed to create a new Product.
type NewProduct struct {
	Title            string  `json:"title"  validate:"required"`
	Slug             string  `json:"slug"  validate:"required"`
	CategoryID       string  `json:"category_id"`
	BrandID          string  `json:"brand_id"`
	Price            float64 `json:"price"`
	OldPrice         float64 `json:"old_price"`
	Image            string  `json:"image"`
	ShortDescription string  `json:"short_description"`
	Description      string  `json:"description"`
	MetaTitle        string  `json:"meta_title"`
	MetaKeywords     string  `json:"meta_keywords"`
	MetaDescription  string  `json:"meta_description"`
}

// UpdateProduct in database
type UpdateProduct struct {
	Title            *string  `json:"title"  validate:"required"`
	Slug             *string  `json:"slug"  validate:"required"`
	CategoryID       *string  `json:"category_id"`
	BrandID          *string  `json:"brand_id"`
	Price            *float64 `json:"price"`
	OldPrice         *float64 `json:"old_price"`
	Image            *string  `json:"image"`
	ShortDescription *string  `json:"short_description"`
	Description      *string  `json:"description"`
	MetaTitle        *string  `json:"meta_title"`
	MetaKeywords     *string  `json:"meta_keywords"`
	MetaDescription  *string  `json:"meta_description"`
}
