package brand_test

import (
	"context"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-cmp/cmp"
	"github.com/igorbelousov/shop-backend/internal/auth"
	"github.com/igorbelousov/shop-backend/internal/data/brand"
	"github.com/igorbelousov/shop-backend/internal/tests"
	"github.com/pkg/errors"
)

func TestBrand(t *testing.T) {
	log, db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)
	testID := 0
	c := brand.New(log, db)
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

	nb := brand.NewBrand{
		Title:           "Test Brand",
		Slug:            "test-brand",
		Image:           "link-to-image",
		Description:     "DESCRIPTION",
		MetaTitle:       "Meta Test brand",
		MetaKeywords:    "meta keywords",
		MetaDescription: "meta description",
	}

	br, err := c.Create(ctx, traceID, claims, nb, now)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to create bregory : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to create bregory.", tests.Success, testID)

	saved, err := c.QueryByID(ctx, traceID, br.ID)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve bregory by ID: %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to retrieve bregory by ID.", tests.Success, testID)

	if diff := cmp.Diff(br, saved); diff != "" {
		t.Fatalf("\t%s\tTest %d:\tShould get back the same bregory. Diff:\n%s", tests.Failed, testID, diff)
	}
	t.Logf("\t%s\tTest %d:\tShould get back the same bregory.", tests.Success, testID)

	upd := brand.UpdateBrand{
		Title: tests.StringPointer("Test bregory Update"),
		Slug:  tests.StringPointer("test-bregory-update"),
	}

	if err := c.Update(ctx, traceID, claims, br.ID, upd, now); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to update brand : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to update brand.", tests.Success, testID)

	saved, err = c.QueryBySlug(ctx, traceID, *upd.Slug)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve bregory by Slug : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to retrieve bregory by Slug.", tests.Success, testID)

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

	if err := c.Delete(ctx, traceID, claims, br.ID); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to delete brand : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to delete brand.", tests.Success, testID)

	_, err = c.QueryByID(ctx, traceID, br.ID)
	if errors.Cause(err) != brand.ErrNotFound {
		t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve brand : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve brand.", tests.Success, testID)

}
