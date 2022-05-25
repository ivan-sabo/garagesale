package handlers

import (
	"log"
	"net/http"

	"github.com/ivan-sabo/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

func API(l *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(l)

	c := Check{db: db}
	app.Handle(http.MethodGet, "/v1/health", c.Health)

	p := Product{DB: db, Log: l}

	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodPost, "/v1/products", p.Create)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)
	app.Handle(http.MethodPut, "/v1/products/{id}", p.Update)
	app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete)

	app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale)
	app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales)

	return app
}
