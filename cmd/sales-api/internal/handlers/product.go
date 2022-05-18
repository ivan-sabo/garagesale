package handlers

import (
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
func (p *Product) List(w http.ResponseWriter, r *http.Request) {
	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("Error querying DB : ", err)
		return
	}

	if err := web.Respond(w, list, http.StatusOK); err != nil {
		p.Log.Println("error responding : ", err)
		return
	}
}

// Retrieve gives a signle Product
func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(p.DB, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("Error querying DB : ", err)
		return
	}

	if err := web.Respond(w, prod, http.StatusOK); err != nil {
		p.Log.Println("error responding : ", err)
		return
	}
}

// Create decode a JSON document from a POST request and create a new Product
func (p *Product) Create(w http.ResponseWriter, r *http.Request) {
	var np product.NewProduct

	if err := web.Decode(r, &np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.Log.Println(err)
		return
	}

	prod, err := product.Create(p.DB, np, time.Now())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error querying DB : ", err)
		return
	}

	if err := web.Respond(w, prod, http.StatusCreated); err != nil {
		p.Log.Println("error responding : ", err)
		return
	}
}
