package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/ivan-sabo/garagesale/internal/platform/web"
	"github.com/ivan-sabo/garagesale/internal/product"
	"github.com/jmoiron/sqlx"
)

// ProductService has handler methods for dealing with Products
type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// ListProducts tells you about request you made
func (p *Product) List(w http.ResponseWriter, r *http.Request) error {
	list, err := product.List(r.Context(), p.DB)
	if err != nil {
		return err
	}

	return web.Respond(w, list, http.StatusOK)
}

// Retrieve gives a signle Product
func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(r.Context(), p.DB, id)
	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("looking for product %q; %w", id, err)
		}
	}

	return web.Respond(w, prod, http.StatusOK)
}

// Create decode a JSON document from a POST request and create a new Product
func (p *Product) Create(w http.ResponseWriter, r *http.Request) error {
	var np product.NewProduct

	if err := web.Decode(r, &np); err != nil {
		return err
	}

	prod, err := product.Create(r.Context(), p.DB, np, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(w, prod, http.StatusCreated)
}
