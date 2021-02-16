package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/igorbelousov/shop-backend/foundation/web"
	"github.com/igorbelousov/shop-backend/internal/auth"
	"github.com/igorbelousov/shop-backend/internal/data/acategory"
	"github.com/igorbelousov/shop-backend/internal/data/article"
	"github.com/igorbelousov/shop-backend/internal/data/brand"
	"github.com/igorbelousov/shop-backend/internal/data/category"
	"github.com/igorbelousov/shop-backend/internal/data/product"
	"github.com/igorbelousov/shop-backend/internal/data/slide"
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

	slide := slideGroup{
		slide: slide.New(log, db),
	}

	brand := brandGroup{
		brand: brand.New(log, db),
	}

	acat := acategoryGroup{
		acategory: acategory.New(log, db),
	}

	art := articleGroup{
		article: article.New(log, db),
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

	app.Handle(http.MethodGet, "/slide/", slide.query)
	app.Handle(http.MethodGet, "/slide/:id", slide.queryByID)
	app.Handle(http.MethodPost, "/slide", slide.create, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodPut, "/slide/:id", slide.update, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/slide/:id", slide.delete, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))

	app.Handle(http.MethodGet, "/brand/", brand.query)
	app.Handle(http.MethodGet, "/brand/:id", brand.queryByID)
	// app.Handle(http.MethodGet, "/brand/:slug", brand.queryBySlug)
	app.Handle(http.MethodPost, "/brand", brand.create, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodPut, "/brand/:id", brand.update, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/brand/:id", brand.delete, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))

	app.Handle(http.MethodGet, "/blog/", acat.query)
	app.Handle(http.MethodGet, "/blog/:id", acat.queryByID)
	// app.Handle(http.MethodGet, "/blog/:slug", acat.queryBySlug)
	app.Handle(http.MethodPost, "/blog", acat.create, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodPut, "/blog/:id", acat.update, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/blog/:id", acat.delete, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))

	app.Handle(http.MethodGet, "/article/", art.query)
	app.Handle(http.MethodGet, "/article/:id", art.queryByID)
	// app.Handle(http.MethodGet, "/article/:slug", art.queryBySlug)
	app.Handle(http.MethodPost, "/article", art.create, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodPut, "/article/:id", art.update, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))
	app.Handle(http.MethodDelete, "/article/:id", art.delete, mid.Authenticate(a), mid.Authorize(auth.RoleAdmin))

	return app
}
