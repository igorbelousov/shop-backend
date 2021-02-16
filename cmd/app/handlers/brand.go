package handlers

import (
	"context"
	"net/http"

	"github.com/igorbelousov/shop-backend/foundation/web"
	"github.com/igorbelousov/shop-backend/internal/auth"
	"github.com/igorbelousov/shop-backend/internal/data/brand"
	"github.com/pkg/errors"
)

type brandGroup struct {
	brand brand.Brand
}

func (bg brandGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	br, err := bg.brand.Query(ctx, v.TraceID)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, br, http.StatusOK)
}

func (bg brandGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	br, err := bg.brand.QueryByID(ctx, v.TraceID, params["id"])
	if err != nil {
		switch err {
		case brand.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case brand.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.Respond(ctx, w, br, http.StatusOK)
}

func (bg brandGroup) queryBySlug(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	br, err := bg.brand.QueryBySlug(ctx, v.TraceID, params["slug"])
	if err != nil {
		switch err {
		case brand.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case brand.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "Slug: %s", params["slug"])
		}
	}

	return web.Respond(ctx, w, br, http.StatusOK)
}

func (bg brandGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("claims missing from context")
	}

	var nb brand.NewBrand
	if err := web.Decode(r, &nb); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	br, err := bg.brand.Create(ctx, v.TraceID, claims, nb, v.Now)
	if err != nil {
		return errors.Wrapf(err, "creating new brand: %+v", nb)
	}

	return web.Respond(ctx, w, br, http.StatusCreated)
}

func (bg brandGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("claims missing from context")
	}

	var upd brand.UpdateBrand
	if err := web.Decode(r, &upd); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	params := web.Params(r)
	if err := bg.brand.Update(ctx, v.TraceID, claims, params["id"], upd, v.Now); err != nil {
		switch err {
		case brand.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case brand.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case brand.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s  User: %+v", params["id"], &upd)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (bg brandGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	if err := bg.brand.Delete(ctx, v.TraceID, claims, params["id"]); err != nil {
		switch err {
		case brand.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
