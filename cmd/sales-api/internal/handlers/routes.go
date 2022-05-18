package handlers

import (
	"log"
	"net/http"

	"github.com/ivan-sabo/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

func API(l *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(l)

	p := Product{DB: db, Log: l}

	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodPost, "/v1/products", p.Create)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)

	return app
}
