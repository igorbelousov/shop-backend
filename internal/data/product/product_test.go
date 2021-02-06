package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-cmp/cmp"
	"github.com/igorbelousov/shop-backend/internal/auth"
	"github.com/igorbelousov/shop-backend/internal/data/product"
	"github.com/igorbelousov/shop-backend/internal/tests"
	"github.com/pkg/errors"
)

func TestProduct(t *testing.T) {
	log, db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)
	testID := 0
	p := product.New(log, db)
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

	np := product.NewProduct{
		Title:       "Test Product",
		Slug:        "test-product",
		CategoryID:  "00000000-0000-0000-0000-000000000000",
		Price:       3123.33,
		Description: "DESCRIPTION",
	}

	prod, err := p.Create(ctx, traceID, claims, np, now)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to create product : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to create product.", tests.Success, testID)

	saved, err := p.QueryByID(ctx, traceID, prod.ID)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve product by ID: %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to retrieve product by ID.", tests.Success, testID)

	if diff := cmp.Diff(prod, saved); diff != "" {
		t.Fatalf("\t%s\tTest %d:\tShould get back the same product. Diff:\n%s", tests.Failed, testID, diff)
	}
	t.Logf("\t%s\tTest %d:\tShould get back the same product.", tests.Success, testID)

	upd := product.UpdateProduct{
		Title: tests.StringPointer("Test Product Update"),
		Slug:  tests.StringPointer("test-product-update"),
	}

	if err := p.Update(ctx, traceID, claims, prod.ID, upd, now); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to update product : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to update product.", tests.Success, testID)

	saved, err = p.QueryBySlug(ctx, traceID, *upd.Slug)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve product by Slug : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to retrieve product by Slug.", tests.Success, testID)

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

	if err := p.Delete(ctx, traceID, claims, prod.ID); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to delete product : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to delete product.", tests.Success, testID)

	_, err = p.QueryByID(ctx, traceID, prod.ID)
	if errors.Cause(err) != product.ErrNotFound {
		t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve product : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve product.", tests.Success, testID)

}
