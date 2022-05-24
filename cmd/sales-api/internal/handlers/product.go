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

// Product defines all of the handlers related to products. It holds the
// application state needed by the handler methods
type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List gets all products from the service layer
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

// Update decodes the body of a request to update an existing product. The ID
// of the product is part of the request URL
func (p *Product) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return fmt.Errorf("decoding product update: %w", err)
	}

	if err := product.Update(r.Context(), p.DB, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("updating product (id: %q): %w", id, err)
		}
	}

	return web.Respond(w, nil, http.StatusNoContent)
}

// Delete removes a single product identified by an ID in the request URL.
func (p *Product) Delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := product.Delete(r.Context(), p.DB, id); err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("deleting product (id: %s): %w", id, err)
		}
	}

	return web.Respond(w, nil, http.StatusNoContent)
}

// AddSale creates a new Sale for a particular product. It looks for a JSON
// object in the request body
func (p *Product) AddSale(w http.ResponseWriter, r *http.Request) error {
	var ns product.NewSale
	if err := web.Decode(r, &ns); err != nil {
		return fmt.Errorf("decoding new sale: %w", err)
	}

	productID := chi.URLParam(r, "id")

	sale, err := product.AddSale(r.Context(), p.DB, ns, productID, time.Now())
	if err != nil {
		return fmt.Errorf("adding new sale: %w", err)
	}

	return web.Respond(w, sale, http.StatusCreated)
}

// ListSales gets all sales for a particular product
func (p *Product) ListSales(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := product.ListSales(r.Context(), p.DB, id)
	if err != nil {
		return fmt.Errorf("getting sales list: %w", err)
	}

	return web.Respond(w, list, http.StatusOK)
}
