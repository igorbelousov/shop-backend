package slide_test

import (
	"context"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-cmp/cmp"
	"github.com/igorbelousov/shop-backend/internal/auth"
	slider "github.com/igorbelousov/shop-backend/internal/data/slide"
	"github.com/igorbelousov/shop-backend/internal/tests"
	"github.com/pkg/errors"
)

func TestSlide(t *testing.T) {
	log, db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)
	testID := 0
	s := slider.New(log, db)
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

	ns := slider.NewSlide{
		Title:    "Test Slide",
		Image:    "test-slide",
		SubTitle: "Sub Title",
		Link:     "http://google.com",
	}

	slide, err := s.Create(ctx, traceID, claims, ns, now)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to create slide : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to create slide.", tests.Success, testID)

	saved, err := s.QueryByID(ctx, traceID, slide.ID)
	if err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve slide by ID: %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to retrieve slide by ID.", tests.Success, testID)

	if diff := cmp.Diff(slide, saved); diff != "" {
		t.Fatalf("\t%s\tTest %d:\tShould get back the same slide. Diff:\n%s", tests.Failed, testID, diff)
	}
	t.Logf("\t%s\tTest %d:\tShould get back the same slide.", tests.Success, testID)

	upd := slider.UpdateSlide{
		Title: tests.StringPointer("Test Slide Update"),
		Image: tests.StringPointer("test-slide-update"),
	}

	if err := s.Update(ctx, traceID, claims, slide.ID, upd, now); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to update slide : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to update slide.", tests.Success, testID)

	if err := s.Delete(ctx, traceID, claims, slide.ID); err != nil {
		t.Fatalf("\t%s\tTest %d:\tShould be able to delete slide : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould be able to delete slide.", tests.Success, testID)

	_, err = s.QueryByID(ctx, traceID, slide.ID)
	if errors.Cause(err) != slider.ErrNotFound {
		t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve slide : %s.", tests.Failed, testID, err)
	}
	t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve slide.", tests.Success, testID)

}
