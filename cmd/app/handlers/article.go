package handlers

import (
	"context"
	"net/http"

	"github.com/igorbelousov/shop-backend/foundation/web"
	"github.com/igorbelousov/shop-backend/internal/auth"
	"github.com/igorbelousov/shop-backend/internal/data/article"
	"github.com/pkg/errors"
)

type articleGroup struct {
	article article.Article
}

func (ag articleGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	ar, err := ag.article.Query(ctx, v.TraceID)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, ar, http.StatusOK)
}

func (ag articleGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	ar, err := ag.article.QueryByID(ctx, v.TraceID, params["id"])
	if err != nil {
		switch err {
		case article.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case article.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.Respond(ctx, w, ar, http.StatusOK)
}

func (ag articleGroup) queryBySlug(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	ar, err := ag.article.QueryBySlug(ctx, v.TraceID, params["slug"])
	if err != nil {
		switch err {
		case article.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case article.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "Slug: %s", params["slug"])
		}
	}

	return web.Respond(ctx, w, ar, http.StatusOK)
}

func (ag articleGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("claims missing from context")
	}

	var na article.NewArticle
	if err := web.Decode(r, &na); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	ar, err := ag.article.Create(ctx, v.TraceID, claims, na, v.Now)
	if err != nil {
		return errors.Wrapf(err, "creating new brand: %+v", na)
	}

	return web.Respond(ctx, w, ar, http.StatusCreated)
}

func (ag articleGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("claims missing from context")
	}

	var upd article.UpdateArticle
	if err := web.Decode(r, &upd); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	params := web.Params(r)
	if err := ag.article.Update(ctx, v.TraceID, claims, params["id"], upd, v.Now); err != nil {
		switch err {
		case article.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case article.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case article.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s  User: %+v", params["id"], &upd)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (ag articleGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	if err := ag.article.Delete(ctx, v.TraceID, claims, params["id"]); err != nil {
		switch err {
		case article.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
