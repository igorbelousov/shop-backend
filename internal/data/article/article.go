package article

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
	// ErrNotFound is used when a specific article is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a article tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

// Article manages the set of API's for article access.
type Article struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a article for api access.
func New(log *log.Logger, db *sqlx.DB) Article {
	return Article{
		log: log,
		db:  db,
	}
}

// Create inserts a new article into the database.
func (a Article) Create(ctx context.Context, traceID string, claims auth.Claims, na NewArticle, now time.Time) (Info, error) {

	// If you are not an admin and looking to retrieve someone elses product.
	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, ErrForbidden
	}

	art := Info{
		ID:              uuid.New().String(),
		Title:           na.Title,
		Slug:            na.Slug,
		CategoryID:      na.CategoryID,
		Image:           na.Image,
		Description:     na.Description,
		MetaTitle:       na.MetaTitle,
		MetaKeywords:    na.MetaKeywords,
		MetaDescription: na.MetaDescription,
		DateCreated:     now.UTC(),
		DateUpdated:     now.UTC(),
	}

	const q = `
	INSERT INTO articles
		(article_id, title, slug, category_id, image, description, meta_title, meta_keywords, meta_description, date_created, date_updated)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	a.log.Printf("%s: %s: %s", traceID, "article.Create",
		database.Log(q, art.ID, art.Title, art.Slug, art.CategoryID, art.Image, art.Description, art.MetaTitle, art.MetaKeywords, art.MetaDescription, art.DateCreated, art.DateUpdated),
	)

	if _, err := a.db.ExecContext(ctx, q, art.ID, art.Title, art.Slug, art.CategoryID, art.Image, art.Description, art.MetaTitle, art.MetaKeywords, art.MetaDescription, art.DateCreated, art.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting article")
	}

	return art, nil
}

// Update replaces a article document in the database.
func (a Article) Update(ctx context.Context, traceID string, claims auth.Claims, articleID string, ua UpdateArticle, now time.Time) error {

	art, err := a.QueryByID(ctx, traceID, articleID)
	if err != nil {
		return err
	}

	if !claims.Authorized(auth.RoleAdmin) {
		return ErrForbidden
	}

	if ua.Title != nil {
		art.Title = *ua.Title
	}
	if ua.Slug != nil {
		art.Slug = *ua.Slug
	}
	if ua.CategoryID != nil {
		art.CategoryID = *ua.CategoryID
	}
	if ua.Image != nil {
		art.Image = *ua.Image
	}
	if ua.Description != nil {
		art.Description = *ua.Description
	}
	if ua.MetaTitle != nil {
		art.MetaTitle = *ua.MetaTitle
	}
	if ua.MetaKeywords != nil {
		art.MetaKeywords = *ua.MetaKeywords
	}
	if ua.MetaDescription != nil {
		art.MetaDescription = *ua.MetaDescription
	}
	art.DateUpdated = now

	const q = `
	UPDATE
		articles
	SET 
		"title" = $2,
		"slug" = $3,
		"category_id" = $4,
		"image" = $5,
		"description" = $6,
		"meta_title" = $7, 
		"meta_keywords" = $8, 
		"meta_description" = $9,
		"date_updated" = $10
	WHERE
		article_id = $1`

	a.log.Printf("%s: %s: %s", traceID, "article.Update",
		database.Log(q, art.ID, art.Title, art.Slug, art.CategoryID, art.Image, art.Description, art.MetaTitle, art.MetaKeywords, art.MetaDescription, art.DateUpdated),
	)

	if _, err = a.db.ExecContext(ctx, q, articleID, art.Title, art.Slug, art.CategoryID, art.Image, art.Description, art.MetaTitle, art.MetaKeywords, art.MetaDescription, art.DateUpdated); err != nil {
		return errors.Wrap(err, "updating article")
	}

	return nil
}

// Delete removes a article from the database.
func (a Article) Delete(ctx context.Context, traceID string, claims auth.Claims, articleID string) error {
	if !claims.Authorized(auth.RoleAdmin) {
		return ErrForbidden
	}
	if _, err := uuid.Parse(articleID); err != nil {
		return ErrInvalidID
	}

	const q = `
	DELETE FROM
		articles
	WHERE
		article_id = $1`

	a.log.Printf("%s: %s: %s", traceID, "article.Delete",
		database.Log(q, articleID),
	)

	if _, err := a.db.ExecContext(ctx, q, articleID); err != nil {
		return errors.Wrapf(err, "deleting article %s", articleID)
	}

	return nil
}

// QueryByID gets the specified article from the database.
func (a Article) QueryByID(ctx context.Context, traceID string, articleID string) (Info, error) {

	if _, err := uuid.Parse(articleID); err != nil {
		return Info{}, ErrInvalidID
	}

	const q = `
	SELECT
		*
	FROM
		articles
	WHERE 
		article_id = $1`

	a.log.Printf("%s: %s: %s", traceID, "article.QueryByID",
		database.Log(q, articleID),
	)

	var cat Info
	if err := a.db.GetContext(ctx, &cat, q, articleID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting article %q", articleID)
	}

	return cat, nil
}

// QueryBySlug gets the specified article from the database.
func (a Article) QueryBySlug(ctx context.Context, traceID string, Slug string) (Info, error) {

	const q = `
	SELECT
		*
	FROM
		articles
	WHERE 
		slug = $1`

	a.log.Printf("%s: %s: %s", traceID, "article.QueryBySlug",
		database.Log(q, Slug),
	)

	var cat Info
	if err := a.db.GetContext(ctx, &cat, q, Slug); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting article %q", Slug)
	}

	return cat, nil
}

// Query retrieves a list of existing article from the database.
func (a Article) Query(ctx context.Context, traceID string) ([]Info, error) {

	const q = `
	SELECT
		*
	FROM
		articles
	ORDER BY
		title`

	a.log.Printf("%s: %s: %s", traceID, "articles.Query",
		database.Log(q),
	)

	categories := []Info{}
	if err := a.db.SelectContext(ctx, &categories, q); err != nil {
		return nil, errors.Wrap(err, "selecting articles")
	}

	return categories, nil
}
