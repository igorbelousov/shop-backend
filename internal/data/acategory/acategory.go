package acategory

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/igorbelousov/shop-backend/foundation/database"
	"github.com/igorbelousov/shop-backend/internal/auth"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is used when a specific arcicle category is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a  arcicle category tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

//  ACategory manages the set of API's for article category access.
type ACategory struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a  arcicle category for api access.
func New(log *log.Logger, db *sqlx.DB) ACategory {
	return ACategory{
		log: log,
		db:  db,
	}
}

// Create inserts a new  arcicle category into the database.
func (c ACategory) Create(ctx context.Context, traceID string, claims auth.Claims, nc NewCategory, now time.Time) (Info, error) {

	// If you are not an admin and looking to retrieve someone elses product.
	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, ErrForbidden
	}

	cat := Info{
		ID:              uuid.New().String(),
		Title:           nc.Title,
		Slug:            nc.Slug,
		Image:           nc.Image,
		Description:     nc.Description,
		MetaTitle:       nc.MetaTitle,
		MetaKeywords:    nc.MetaKeywords,
		MetaDescription: nc.MetaDescription,
		DateCreated:     now.UTC(),
		DateUpdated:     now.UTC(),
	}

	const q = `
	INSERT INTO article_categories
		(category_id, title, slug,  image, description, meta_title, meta_keywords, meta_description, date_created, date_updated)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	c.log.Printf("%s: %s: %s", traceID, "article-category.Create",
		database.Log(q, cat.ID, cat.Title, cat.Slug, cat.Image, cat.Description, cat.MetaTitle, cat.MetaKeywords, cat.MetaDescription, cat.DateCreated, cat.DateUpdated),
	)

	if _, err := c.db.ExecContext(ctx, q, cat.ID, cat.Title, cat.Slug, cat.Image, cat.Description, cat.MetaTitle, cat.MetaKeywords, cat.MetaDescription, cat.DateCreated, cat.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting article-category")
	}

	return cat, nil
}

// Update replaces a  arcicle category document in the database.
func (c ACategory) Update(ctx context.Context, traceID string, claims auth.Claims, categoryID string, uc UpdateCategory, now time.Time) error {

	cat, err := c.QueryByID(ctx, traceID, categoryID)
	if err != nil {
		return err
	}

	if !claims.Authorized(auth.RoleAdmin) {
		return ErrForbidden
	}

	if uc.Title != nil {
		cat.Title = *uc.Title
	}
	if uc.Slug != nil {
		cat.Slug = *uc.Slug
	}
	if uc.Image != nil {
		cat.Image = *uc.Image
	}
	if uc.Description != nil {
		cat.Description = *uc.Description
	}
	if uc.MetaTitle != nil {
		cat.MetaTitle = *uc.MetaTitle
	}
	if uc.MetaKeywords != nil {
		cat.MetaKeywords = *uc.MetaKeywords
	}
	if uc.MetaDescription != nil {
		cat.MetaDescription = *uc.MetaDescription
	}
	cat.DateUpdated = now

	const q = `
	UPDATE
		article_categories
	SET 
		"title" = $2,
		"slug" = $3,
		"image" = $4,
		"description" = $5,
		"meta_title" = $6, 
		"meta_keywords" = $7, 
		"meta_description" = $8,
		"date_updated" = $9
	WHERE
		category_id = $1`

	c.log.Printf("%s: %s: %s", traceID, "carticle-ategory.Update",
		database.Log(q, cat.ID, cat.Title, cat.Slug, cat.Image, cat.Description, cat.MetaTitle, cat.MetaKeywords, cat.MetaDescription, cat.DateUpdated),
	)

	if _, err = c.db.ExecContext(ctx, q, categoryID, cat.Title, cat.Slug, cat.Image, cat.Description, cat.MetaTitle, cat.MetaKeywords, cat.MetaDescription, cat.DateUpdated); err != nil {
		return errors.Wrap(err, "updating article-category")
	}

	return nil
}

// Delete removes a  arcicle category from the database.
func (c ACategory) Delete(ctx context.Context, traceID string, claims auth.Claims, categoryID string) error {
	if !claims.Authorized(auth.RoleAdmin) {
		return ErrForbidden
	}
	if _, err := uuid.Parse(categoryID); err != nil {
		return ErrInvalidID
	}

	const q = `
	DELETE FROM
		article_categories
	WHERE
		category_id = $1`

	c.log.Printf("%s: %s: %s", traceID, "article-category.Delete",
		database.Log(q, categoryID),
	)

	if _, err := c.db.ExecContext(ctx, q, categoryID); err != nil {
		return errors.Wrapf(err, "deleting article-category %s", categoryID)
	}

	return nil
}

// QueryByID gets the specified  arcicle category from the database.
func (c ACategory) QueryByID(ctx context.Context, traceID string, categoryID string) (Info, error) {

	if _, err := uuid.Parse(categoryID); err != nil {
		return Info{}, ErrInvalidID
	}

	const q = `
	SELECT
		*
	FROM
		article_categories
	WHERE 
		category_id = $1`

	c.log.Printf("%s: %s: %s", traceID, "article-category.QueryByID",
		database.Log(q, categoryID),
	)

	var cat Info
	if err := c.db.GetContext(ctx, &cat, q, categoryID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting article-category %q", categoryID)
	}

	return cat, nil
}

// QueryBySlug gets the specified  arcicle category from the database.
func (c ACategory) QueryBySlug(ctx context.Context, traceID string, Slug string) (Info, error) {

	const q = `
	SELECT
		*
	FROM
		article_categories
	WHERE 
		slug = $1`

	c.log.Printf("%s: %s: %s", traceID, "article-category.QueryBySlug",
		database.Log(q, Slug),
	)

	var cat Info
	if err := c.db.GetContext(ctx, &cat, q, Slug); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting article-category %q", Slug)
	}

	return cat, nil
}

// Query retrieves a list of existing  arcicle category from the database.
func (c ACategory) Query(ctx context.Context, traceID string) ([]Info, error) {

	const q = `
	SELECT
		*
	FROM
		article_categories
	ORDER BY
		title`

	c.log.Printf("%s: %s: %s", traceID, "article-categories.Query",
		database.Log(q),
	)

	categories := []Info{}
	if err := c.db.SelectContext(ctx, &categories, q); err != nil {
		return nil, errors.Wrap(err, "selecting article-categories")
	}

	return categories, nil
}
