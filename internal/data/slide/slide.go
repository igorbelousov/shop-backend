package slide

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
	// ErrNotFound is used when a specific Slide is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a slide tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

// Slide manages the set of API's for slide access.
type Slide struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a Slide for api access.
func New(log *log.Logger, db *sqlx.DB) Slide {
	return Slide{
		log: log,
		db:  db,
	}
}

// Create inserts a new slide into the database.
func (s Slide) Create(ctx context.Context, traceID string, claims auth.Claims, ns NewSlide, now time.Time) (Info, error) {

	// If you are not an admin and looking to retrieve someone elses slide.
	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, ErrForbidden
	}

	slide := Info{
		ID:          uuid.New().String(),
		Title:       ns.Title,
		Image:       ns.Image,
		SubTitle:    ns.SubTitle,
		Link:        ns.Link,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `
	INSERT INTO slides
		(slide_id, title, image, sub_title, link, date_created, date_updated)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)`

	s.log.Printf("%s: %s: %s", traceID, "slide.Create",
		database.Log(q, slide.ID, slide.Title, slide.Image, slide.SubTitle, slide.Link, slide.DateCreated, slide.DateUpdated),
	)

	if _, err := s.db.ExecContext(ctx, q, slide.ID, slide.Title, slide.Image, slide.SubTitle, slide.Link, slide.DateCreated, slide.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting slide")
	}

	return slide, nil
}

// Update replaces a slide document in the database.
func (s Slide) Update(ctx context.Context, traceID string, claims auth.Claims, SlideID string, us UpdateSlide, now time.Time) error {

	slide, err := s.QueryByID(ctx, traceID, SlideID)
	if err != nil {
		return err
	}

	if !claims.Authorized(auth.RoleAdmin) {
		return ErrForbidden
	}

	if us.Title != nil {
		slide.Title = *us.Title
	}
	if us.Image != nil {
		slide.Image = *us.Image
	}
	if us.SubTitle != nil {
		slide.SubTitle = *us.SubTitle
	}
	if us.Link != nil {
		slide.Link = *us.Link
	}
	slide.DateUpdated = now

	const q = `
	UPDATE
		slides
	SET 
		"title" = $2,
		"image" = $3,
		"sub_title" = $4,
		"link" = $5,
		"date_updated" = $6
	WHERE
		slide_id = $1`

	s.log.Printf("%s: %s: %s", traceID, "slide.Update",
		database.Log(q, slide.ID, slide.Title, slide.Image, slide.SubTitle, slide.Link, slide.DateUpdated),
	)

	if _, err = s.db.ExecContext(ctx, q, SlideID, slide.Title, slide.Image, slide.SubTitle, slide.Link, slide.DateUpdated); err != nil {
		return errors.Wrap(err, "updating slide")
	}

	return nil
}

// Delete removes a slide from the database.
func (s Slide) Delete(ctx context.Context, traceID string, claims auth.Claims, SlideID string) error {
	if !claims.Authorized(auth.RoleAdmin) {
		return ErrForbidden
	}
	if _, err := uuid.Parse(SlideID); err != nil {
		return ErrInvalidID
	}

	const q = `
	DELETE FROM
		slides
	WHERE
		slide_id = $1`

	s.log.Printf("%s: %s: %s", traceID, "slide.Delete",
		database.Log(q, SlideID),
	)

	if _, err := s.db.ExecContext(ctx, q, SlideID); err != nil {
		return errors.Wrapf(err, "deleting cateslidegory %s", SlideID)
	}

	return nil
}

// QueryByID gets the specified slide from the database.
func (s Slide) QueryByID(ctx context.Context, traceID string, SlideID string) (Info, error) {

	if _, err := uuid.Parse(SlideID); err != nil {
		return Info{}, ErrInvalidID
	}

	const q = `
	SELECT
		*
	FROM
		slides
	WHERE 
		slide_id = $1`

	s.log.Printf("%s: %s: %s", traceID, "slide.QueryByID",
		database.Log(q, SlideID),
	)

	var slide Info
	if err := s.db.GetContext(ctx, &slide, q, SlideID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting slide %q", SlideID)
	}

	return slide, nil
}

// Query retrieves a list of existing slides from the database.
func (s Slide) Query(ctx context.Context, traceID string) ([]Info, error) {

	const q = `
	SELECT
		*
	FROM
		slides
	ORDER BY
		date_created`

	s.log.Printf("%s: %s: %s", traceID, "slides.Query",
		database.Log(q),
	)

	slides := []Info{}
	if err := s.db.SelectContext(ctx, &slides, q); err != nil {
		return nil, errors.Wrap(err, "selecting slides")
	}

	return slides, nil
}
