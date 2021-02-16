package brand

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
	// ErrNotFound is used when a specific brand is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

// brand manages the set of API's for user access.
type Brand struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a brand for api access.
func New(log *log.Logger, db *sqlx.DB) Brand {
	return Brand{
		log: log,
		db:  db,
	}
}

// Create inserts a new brand into the database.
func (b Brand) Create(ctx context.Context, traceID string, claims auth.Claims, nb NewBrand, now time.Time) (Info, error) {

	// If you are not an admin and looking to retrieve someone elses brand.
	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, ErrForbidden
	}

	br := Info{
		ID:              uuid.New().String(),
		Title:           nb.Title,
		Slug:            nb.Slug,
		Image:           nb.Image,
		Description:     nb.Description,
		MetaTitle:       nb.MetaTitle,
		MetaKeywords:    nb.MetaKeywords,
		MetaDescription: nb.MetaDescription,
		DateCreated:     now.UTC(),
		DateUpdated:     now.UTC(),
	}

	const q = `
	INSERT INTO brands
		(brand_id, title, slug,  image, description, meta_title, meta_keywords, meta_description, date_created, date_updated)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	b.log.Printf("%s: %s: %s", traceID, "brand.Create",
		database.Log(q, br.ID, br.Title, br.Slug, br.Image, br.Description, br.MetaTitle, br.MetaKeywords, br.MetaDescription, br.DateCreated, br.DateUpdated),
	)

	if _, err := b.db.ExecContext(ctx, q, br.ID, br.Title, br.Slug, br.Image, br.Description, br.MetaTitle, br.MetaKeywords, br.MetaDescription, br.DateCreated, br.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting brand")
	}

	return br, nil
}

// Update replaces a brand document in the database.
func (b Brand) Update(ctx context.Context, traceID string, claims auth.Claims, brandID string, ub UpdateBrand, now time.Time) error {

	br, err := b.QueryByID(ctx, traceID, brandID)
	if err != nil {
		return err
	}

	if !claims.Authorized(auth.RoleAdmin) {
		return ErrForbidden
	}

	if ub.Title != nil {
		br.Title = *ub.Title
	}
	if ub.Slug != nil {
		br.Slug = *ub.Slug
	}
	if ub.Image != nil {
		br.Image = *ub.Image
	}
	if ub.Description != nil {
		br.Description = *ub.Description
	}
	if ub.MetaTitle != nil {
		br.MetaTitle = *ub.MetaTitle
	}
	if ub.MetaKeywords != nil {
		br.MetaKeywords = *ub.MetaKeywords
	}
	if ub.MetaDescription != nil {
		br.MetaDescription = *ub.MetaDescription
	}
	br.DateUpdated = now

	const q = `
	UPDATE
		brands
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
		brand_id = $1`

	b.log.Printf("%s: %s: %s", traceID, "brand.Update",
		database.Log(q, br.ID, br.Title, br.Slug, br.Image, br.Description, br.MetaTitle, br.MetaKeywords, br.MetaDescription, br.DateUpdated),
	)

	if _, err = b.db.ExecContext(ctx, q, brandID, br.Title, br.Slug, br.Image, br.Description, br.MetaTitle, br.MetaKeywords, br.MetaDescription, br.DateUpdated); err != nil {
		return errors.Wrap(err, "updating brand")
	}

	return nil
}

// Delete removes a brand from the database.
func (b Brand) Delete(ctx context.Context, traceID string, claims auth.Claims, brandID string) error {
	if !claims.Authorized(auth.RoleAdmin) {
		return ErrForbidden
	}
	if _, err := uuid.Parse(brandID); err != nil {
		return ErrInvalidID
	}

	const q = `
	DELETE FROM
		brands
	WHERE
		brand_id = $1`

	b.log.Printf("%s: %s: %s", traceID, "brand.Delete",
		database.Log(q, brandID),
	)

	if _, err := b.db.ExecContext(ctx, q, brandID); err != nil {
		return errors.Wrapf(err, "deleting brand %s", brandID)
	}

	return nil
}

// QueryByID gets the specified brand from the database.
func (b Brand) QueryByID(ctx context.Context, traceID string, brandID string) (Info, error) {

	if _, err := uuid.Parse(brandID); err != nil {
		return Info{}, ErrInvalidID
	}

	const q = `
	SELECT
		*
	FROM
		brands
	WHERE 
		brand_id = $1`

	b.log.Printf("%s: %s: %s", traceID, "brand.QueryByID",
		database.Log(q, brandID),
	)

	var br Info
	if err := b.db.GetContext(ctx, &br, q, brandID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting brand %q", brandID)
	}

	return br, nil
}

// QueryBySlug gets the specified brand from the database.
func (b Brand) QueryBySlug(ctx context.Context, traceID string, Slug string) (Info, error) {

	const q = `
	SELECT
		*
	FROM
		brands
	WHERE 
		slug = $1`

	b.log.Printf("%s: %s: %s", traceID, "brand.QueryBySlug",
		database.Log(q, Slug),
	)

	var br Info
	if err := b.db.GetContext(ctx, &br, q, Slug); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting brand %q", Slug)
	}

	return br, nil
}

// Query retrieves a list of existing brand from the database.
func (b Brand) Query(ctx context.Context, traceID string) ([]Info, error) {

	const q = `
	SELECT
		*
	FROM
		categories
	ORDER BY
		title`

	b.log.Printf("%s: %s: %s", traceID, "brand.Query",
		database.Log(q),
	)

	br := []Info{}
	if err := b.db.SelectContext(ctx, &br, q); err != nil {
		return nil, errors.Wrap(err, "selecting brand")
	}

	return br, nil
}
