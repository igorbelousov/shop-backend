package handlers

import (
	"context"
	"net/http"

	"github.com/igorbelousov/shop-backend/foundation/web"
	"github.com/igorbelousov/shop-backend/internal/auth"
	"github.com/igorbelousov/shop-backend/internal/data/slide"
	"github.com/pkg/errors"
)

type slideGroup struct {
	slide slide.Slide
}

func (sg slideGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	sld, err := sg.slide.Query(ctx, v.TraceID)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, sld, http.StatusOK)
}

func (sg slideGroup) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	sld, err := sg.slide.QueryByID(ctx, v.TraceID, params["id"])
	if err != nil {
		switch err {
		case slide.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case slide.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.Respond(ctx, w, sld, http.StatusOK)
}

func (sg slideGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("claims missing from context")
	}

	var ns slide.NewSlide
	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	sld, err := sg.slide.Create(ctx, v.TraceID, claims, ns, v.Now)
	if err != nil {
		return errors.Wrapf(err, "creating new slide: %+v", ns)
	}

	return web.Respond(ctx, w, sld, http.StatusCreated)
}

func (sg slideGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("claims missing from context")
	}

	var upd slide.UpdateSlide
	if err := web.Decode(r, &upd); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	params := web.Params(r)
	if err := sg.slide.Update(ctx, v.TraceID, claims, params["id"], upd, v.Now); err != nil {
		switch err {
		case slide.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case slide.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case slide.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s  User: %+v", params["id"], &upd)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (sg slideGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	params := web.Params(r)
	if err := sg.slide.Delete(ctx, v.TraceID, claims, params["id"]); err != nil {
		switch err {
		case slide.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
