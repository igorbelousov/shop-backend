package article_test

import (
	"context"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-cmp/cmp"
	"github.com/igorbelousov/shop-backend/internal/auth"
	"github.com/igorbelousov/shop-backend/internal/data/article"
	"github.com/igorbelousov/shop-backend/internal/tests"
	"github.com/pkg/errors"
)

func TestCategory(t *testing.T) {
	log, db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)
	testID := 0
	a := article.New(log, db)
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

	na := article.NewArticle{
		Title:           "Test article",
		Slug:            "test-article",
		CategoryID:      "6628112c-f80f-449d-b5d4-95d71b8472b8",
		Image:           "link-to-image",
		Description:     "DESCRIPTION",
		MetaTitle:       "Meta Test article",
		MetaKeywords:    "meta keywords",
		MetaDescription: "meta description",
	}

	art, err := a.Create(ctx, traceID, claims, na, now)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to create article : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to create article.", tests.Success, testID)

	saved, err := a.QueryByID(ctx, traceID, art.ID)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve article by ID: %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to retrieve article by ID.", tests.Success, testID)

	if diff := cmp.Diff(art, saved); diff != "" {
		t.Fatalf("\t%s\tTest %d:\tShould get back the same article. Diff:\n%s", tests.Failed, testID, diff)
	}
	t.Logf("\t%s\tTest %d:\tShould get back the same article.", tests.Success, testID)

	upd := article.UpdateArticle{
		Title: tests.StringPointer("Test Article Update"),
		Slug:  tests.StringPointer("test-article-update"),
	}

	if err := a.Update(ctx, traceID, claims, art.ID, upd, now); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to update article : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to update article.", tests.Success, testID)

	saved, err = a.QueryBySlug(ctx, traceID, *upd.Slug)
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

	if err := a.Delete(ctx, traceID, claims, art.ID); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to delete article : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to delete article.", tests.Success, testID)

	_, err = a.QueryByID(ctx, traceID, art.ID)
	if errors.Cause(err) != article.ErrNotFound {
		t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve article : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve article.", tests.Success, testID)

}
