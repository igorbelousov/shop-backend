package product

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

// Product manages the set of API's for user access.
type Product struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a Category for api access.
func New(log *log.Logger, db *sqlx.DB) Product {
	return Product{
		log: log,
		db:  db,
	}
}

// Create inserts a new user into the database.
func (p Product) Create(ctx context.Context, traceID string, claims auth.Claims, np NewProduct, now time.Time) (Info, error) {

	// If you are not an admin and looking to retrieve someone elses product.
	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, ErrForbidden
	}

	prod := Info{
		ID:          uuid.New().String(),
		Title:       np.Title,
		Slug:        np.Slug,
		CategoryID:  np.CategoryID,
		Price:       np.Price,
		Description: np.Description,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `
	INSERT INTO products
		(product_id, title, slug, category_id, price, description, date_created, date_updated)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)`

	p.log.Printf("%s: %s: %s", traceID, "product.Create",
		database.Log(q, prod.ID, prod.Title, prod.Slug, prod.CategoryID, prod.Slug, prod.Description, prod.DateCreated, prod.DateUpdated),
	)

	if _, err := p.db.ExecContext(ctx, q, prod.ID, prod.Title, prod.Slug, prod.CategoryID, prod.Slug, prod.Description, prod.DateCreated, prod.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting product")
	}

	return prod, nil
}

// Update replaces a product document in the database.
func (p Product) Update(ctx context.Context, traceID string, claims auth.Claims, productID string, up UpdateProduct, now time.Time) error {

	prod, err := p.QueryByID(ctx, traceID, productID)
	if err != nil {
		return err
	}

	if !claims.Authorized(auth.RoleAdmin) {
		return ErrForbidden
	}

	if up.Title != nil {
		prod.Title = up.Title
	}
	if up.Slug != nil {
		prod.Slug = up.Slug
	}
	if up.CategoryID != nil {
		prod.CategoryID = up.CategoryID
	}
	if up.Price != nil {
		prod.Price = up.Price
	}
	if up.Description != nil {
		prod.Description = up.Description
	}
	prod.DateUpdated = now

	const q = `
	UPDATE
		products	
	SET 
		"title" = $2,
		"slug" = $3,
		"category_id" = $4,
		"price" = $5,
		"description" = $6,
		"date_updated" = $7
	WHERE
		product_id = $1`

	p.log.Printf("%s: %s: %s", traceID, "product.Update",
		database.Log(q, prod.ID, prod.Title, prod.Slug, prod.CategoryID, prod.Description, prod.DateCreated, prod.DateUpdated),
	)

	if _, err = p.db.ExecContext(ctx, q, prod.ID, prod.Title, prod.Slug, prod.CategoryID, prod.Description, prod.DateUpdated); err != nil {
		return errors.Wrap(err, "updating product")
	}

	return nil
}

// Delete removes a product from the database.
func (p Product) Delete(ctx context.Context, traceID string, claims auth.Claims, productID string) error {
	if !claims.Authorized(auth.RoleAdmin) {
		return ErrForbidden
	}
	if _, err := uuid.Parse(productID); err != nil {
		return ErrInvalidID
	}

	const q = `
	DELETE FROM
		products
	WHERE
		product_id = $1`

	p.log.Printf("%s: %s: %s", traceID, "product.Delete",
		database.Log(q, productID),
	)

	if _, err := p.db.ExecContext(ctx, q, productID); err != nil {
		return errors.Wrapf(err, "deleting product %s", productID)
	}

	return nil
}

// QueryByID gets the specified product from the database.
func (p Product) QueryByID(ctx context.Context, traceID string, productID string) (Info, error) {

	if _, err := uuid.Parse(productID); err != nil {
		return Info{}, ErrInvalidID
	}

	const q = `
	SELECT
		*
	FROM
		products
	WHERE 
		product_id = $1`

	p.log.Printf("%s: %s: %s", traceID, "product.QueryByID",
		database.Log(q, productID),
	)

	var cat Info
	if err := p.db.GetContext(ctx, &cat, q, productID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting product %q", productID)
	}

	return cat, nil
}

// QueryBySlug gets the specified product from the database.
func (p Product) QueryBySlug(ctx context.Context, traceID string, Slug string) (Info, error) {

	const q = `
	SELECT
		*
	FROM
		product_id
	WHERE 
		slug = $1`

	p.log.Printf("%s: %s: %s", traceID, "product.QueryBySlug",
		database.Log(q, Slug),
	)

	var cat Info
	if err := p.db.GetContext(ctx, &cat, q, Slug); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting product %q", Slug)
	}

	return cat, nil
}

// Query retrieves a list of existing product from the database.
func (p Product) Query(ctx context.Context, traceID string) ([]Info, error) {

	const q = `
	SELECT
		*
	FROM
		products
	ORDER BY
		date_created`

	p.log.Printf("%s: %s: %s", traceID, "products.Query",
		database.Log(q),
	)

	categories := []Info{}
	if err := p.db.SelectContext(ctx, &categories, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return categories, nil
}
