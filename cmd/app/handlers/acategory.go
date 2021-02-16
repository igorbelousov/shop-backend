package handlers

import (
	"context"
	"net/http"

	"github.com/igorbelousov/shop-backend/foundation/web"
	"github.com/igorbelousov/shop-backend/internal/auth"
	"github.com/igorbelousov/shop-backend/internal/data/acategory"
	"github.com/pkg/errors"
)

type acategoryGroup struct {
	acategory acategory.ACategory
}

func (cg acategoryGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	category, err := cg.acategory.Query(ctx, v.TraceID)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, category, http.StatusOK)
}

func (cg acategoryGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	cat, err := cg.acategory.QueryByID(ctx, v.TraceID, params["id"])
	if err != nil {
		switch err {
		case acategory.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case acategory.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.Respond(ctx, w, cat, http.StatusOK)
}

func (cg acategoryGroup) queryBySlug(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	cat, err := cg.acategory.QueryBySlug(ctx, v.TraceID, params["slug"])
	if err != nil {
		switch err {
		case acategory.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case acategory.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "Slug: %s", params["slug"])
		}
	}

	return web.Respond(ctx, w, cat, http.StatusOK)
}

func (cg acategoryGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("claims missing from context")
	}

	var nc acategory.NewCategory
	if err := web.Decode(r, &nc); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	cat, err := cg.acategory.Create(ctx, v.TraceID, claims, nc, v.Now)
	if err != nil {
		return errors.Wrapf(err, "creating new product: %+v", nc)
	}

	return web.Respond(ctx, w, cat, http.StatusCreated)
}

func (cg acategoryGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("claims missing from context")
	}

	var upd acategory.UpdateCategory
	if err := web.Decode(r, &upd); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	params := web.Params(r)
	if err := cg.acategory.Update(ctx, v.TraceID, claims, params["id"], upd, v.Now); err != nil {
		switch err {
		case acategory.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case acategory.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case acategory.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s  User: %+v", params["id"], &upd)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (cg acategoryGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	if err := cg.acategory.Delete(ctx, v.TraceID, claims, params["id"]); err != nil {
		switch err {
		case acategory.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
