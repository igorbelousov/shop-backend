package category

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
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

// Category manages the set of API's for user access.
type Category struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a Category for api access.
func New(log *log.Logger, db *sqlx.DB) Category {
	return Category{
		log: log,
		db:  db,
	}
}

// Create inserts a new user into the database.
func (c Category) Create(ctx context.Context, traceID string, claims auth.Claims, nc NewCategory, now time.Time) (Info, error) {

	// If you are not an admin and looking to retrieve someone elses product.
	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, ErrForbidden
	}

	cat := Info{
		ID:          uuid.New().String(),
		Title:       nc.Title,
		Slug:        nc.Slug,
		ParrentID:   nc.ParrentID,
		Description: nc.Description,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `
	INSERT INTO categories
		(category_id, title, slug, parrent_id, description, date_created, date_updated)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)`

	c.log.Printf("%s: %s: %s", traceID, "category.Create",
		database.Log(q, cat.ID, cat.Title, cat.Slug, cat.ParrentID, cat.Description, cat.DateCreated, cat.DateUpdated),
	)

	if _, err := c.db.ExecContext(ctx, q, cat.ID, cat.Title, cat.Slug, cat.ParrentID, cat.Description, cat.DateCreated, cat.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting category")
	}

	return cat, nil
}

// Update replaces a CATEGORY document in the database.
func (c Category) Update(ctx context.Context, traceID string, claims auth.Claims, categoryID string, uc UpdateCategory, now time.Time) error {

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
	if uc.ParrentID != nil {
		cat.ParrentID = *uc.ParrentID
	}
	if uc.Description != nil {
		cat.Description = *uc.Description
	}
	cat.DateUpdated = now

	const q = `
	UPDATE
		categories
	SET 
		"title" = $2,
		"slug" = $3,
		"parrent_id" = $4,
		"description" = $5,
		"date_updated" = $6
	WHERE
		category_id = $1`

	c.log.Printf("%s: %s: %s", traceID, "category.Update",
		database.Log(q, cat.ID, cat.Title, cat.Slug, cat.ParrentID, cat.Description, cat.DateCreated, cat.DateUpdated),
	)

	if _, err = c.db.ExecContext(ctx, q, categoryID, cat.Title, cat.Slug, cat.ParrentID, cat.Description, cat.DateUpdated); err != nil {
		return errors.Wrap(err, "updating category")
	}

	return nil
}

// Delete removes a category from the database.
func (c Category) Delete(ctx context.Context, traceID string, claims auth.Claims, categoryID string) error {
	if !claims.Authorized(auth.RoleAdmin) {
		return ErrForbidden
	}
	if _, err := uuid.Parse(categoryID); err != nil {
		return ErrInvalidID
	}

	const q = `
	DELETE FROM
		categories
	WHERE
		category_id = $1`

	c.log.Printf("%s: %s: %s", traceID, "category.Delete",
		database.Log(q, categoryID),
	)

	if _, err := c.db.ExecContext(ctx, q, categoryID); err != nil {
		return errors.Wrapf(err, "deleting category %s", categoryID)
	}

	return nil
}

// QueryByID gets the specified category from the database.
func (c Category) QueryByID(ctx context.Context, traceID string, categoryID string) (Info, error) {

	if _, err := uuid.Parse(categoryID); err != nil {
		return Info{}, ErrInvalidID
	}

	const q = `
	SELECT
		*
	FROM
		categories
	WHERE 
		category_id = $1`

	c.log.Printf("%s: %s: %s", traceID, "category.QueryByID",
		database.Log(q, categoryID),
	)

	var cat Info
	if err := c.db.GetContext(ctx, &cat, q, categoryID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting category %q", categoryID)
	}

	return cat, nil
}

// QueryBySlug gets the specified Category from the database.
func (c Category) QueryBySlug(ctx context.Context, traceID string, Slug string) (Info, error) {

	const q = `
	SELECT
		*
	FROM
		categories
	WHERE 
		slug = $1`

	c.log.Printf("%s: %s: %s", traceID, "category.QueryBySlug",
		database.Log(q, Slug),
	)

	var cat Info
	if err := c.db.GetContext(ctx, &cat, q, Slug); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting category %q", Slug)
	}

	return cat, nil
}

// Query retrieves a list of existing Categories from the database.
func (c Category) Query(ctx context.Context, traceID string) ([]Info, error) {

	const q = `
	SELECT
		*
	FROM
		categories
	ORDER BY
		title`

	c.log.Printf("%s: %s: %s", traceID, "categories.Query",
		database.Log(q),
	)

	categories := []Info{}
	if err := c.db.SelectContext(ctx, &categories, q); err != nil {
		return nil, errors.Wrap(err, "selecting categories")
	}

	return categories, nil
}
