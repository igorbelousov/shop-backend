package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/igorbelousov/shop-backend/foundation/web"
	"github.com/igorbelousov/shop-backend/internal/auth"
	"github.com/igorbelousov/shop-backend/internal/data/category"
	"github.com/igorbelousov/shop-backend/internal/data/product"
	"github.com/igorbelousov/shop-backend/internal/data/user"
	"github.com/igorbelousov/shop-backend/internal/mid"
	"github.com/jmoiron/sqlx"
)

//API function for define routers
func API(build string, shutdown chan os.Signal, log *log.Logger, a *auth.Auth, db *sqlx.DB) *web.App {

	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	cg := checkGroup{
		build: build,
		db:    db,
	}

	app.Handle(http.MethodGet, "/readiness", cg.readiness)
	app.Handle(http.MethodGet, "/liveiness", cg.liveness)

	ug := userGroup{
		user: user.New(log, db),
		auth: a,
	}

	catg := categoryGroup{
		category: category.New(log, db),
	}

	prod := productGroup{
		product: product.New(log, db),
	}

	app.Handle(http.MethodGet, "/users/:page/:rows", ug.query, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodGet, "/users/token/:kid", ug.token)
	app.Handle(http.MethodGet, "/users/:id", ug.queryByID, mid.Authenticate(a))
	app.Handle(http.MethodPost, "/users", ug.create, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodPut, "/users/:id", ug.update, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/users/:id", ug.delete, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))

	app.Handle(http.MethodGet, "/category/", catg.query)
	app.Handle(http.MethodGet, "/category/:id", catg.queryByID)
	// app.Handle(http.MethodGet, "/category/:slug", catg.queryBySlug)
	app.Handle(http.MethodPost, "/category", catg.create, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodPut, "/category/:id", catg.update, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/category/:id", catg.delete, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))

	app.Handle(http.MethodGet, "/product/", prod.query)
	app.Handle(http.MethodGet, "/product/:id", prod.queryByID)
	// app.Handle(http.MethodGet, "/category/:slug", prod.queryBySlug)
	app.Handle(http.MethodPost, "/product", prod.create, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodPut, "/product/:id", prod.update, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/product/:id", prod.delete, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))

	return app
}
