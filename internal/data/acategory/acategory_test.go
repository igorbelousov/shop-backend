package acategory_test

import (
	"context"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-cmp/cmp"
	"github.com/igorbelousov/shop-backend/internal/auth"
	"github.com/igorbelousov/shop-backend/internal/data/acategory"
	"github.com/igorbelousov/shop-backend/internal/tests"
	"github.com/pkg/errors"
)

func TestCategory(t *testing.T) {
	log, db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)
	testID := 0
	c := acategory.New(log, db)
	ctx := context.Background()
	now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)
	traceID := "00000000-0000-0000-0000-000000000000"

	claims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Shop backend",
			Subject:   "00000000-0000-0000-0000-000000000000",
			Audience:  "SHOP",
			ExpiresAt: now.Add(time.Hour).Unix(),
			IssuedAt:  now.Unix(),
		},
		Roles: []string{auth.RoleAdmin},
	}

	nc := acategory.NewCategory{
		Title:           "Test Article Category",
		Slug:            "test-article-category",
		Image:           "link-to-image",
		Description:     "DESCRIPTION",
		MetaTitle:       "Meta Test Category",
		MetaKeywords:    "meta keywords",
		MetaDescription: "meta description",
	}

	cat, err := c.Create(ctx, traceID, claims, nc, now)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to create category : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to create category.", tests.Success, testID)

	saved, err := c.QueryByID(ctx, traceID, cat.ID)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve category by ID: %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to retrieve category by ID.", tests.Success, testID)

	if diff := cmp.Diff(cat, saved); diff != "" {
		t.Fatalf("\t%s\tTest %d:\tShould get back the same category. Diff:\n%s", tests.Failed, testID, diff)
	}
	t.Logf("\t%s\tTest %d:\tShould get back the same category.", tests.Success, testID)

	upd := acategory.UpdateCategory{
		Title: tests.StringPointer("Test Category Update"),
		Slug:  tests.StringPointer("test-category-update"),
	}

	if err := c.Update(ctx, traceID, claims, cat.ID, upd, now); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to update category : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to update category.", tests.Success, testID)

	saved, err = c.QueryBySlug(ctx, traceID, *upd.Slug)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve category by Slug : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to retrieve category by Slug.", tests.Success, testID)

	if saved.Title != *upd.Title {
		t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Title.", tests.Failed, testID)
		t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Title)
		t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Title)
	} else {
		t.Logf("\t%s\tTest %d:\tShould be able to see updates to Title.", tests.Success, testID)
	}

	if saved.Slug != *upd.Slug {
		t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Slug.", tests.Failed, testID)
		t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Slug)
		t.Logf("\t\tTest %d:\tExp: %v", testID, *upd.Slug)
	} else {
		t.Logf("\t%s\tTest %d:\tShould be able to see updates to Slug.", tests.Success, testID)
	}

	if err := c.Delete(ctx, traceID, claims, cat.ID); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to delete category : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to delete category.", tests.Success, testID)

	_, err = c.QueryByID(ctx, traceID, cat.ID)
	if errors.Cause(err) != acategory.ErrNotFound {
		t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve category : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve category.", tests.Success, testID)

}
