package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ivan-sabo/garagesale/internal/product"
	"github.com/jmoiron/sqlx"
)

// ProductService has handler methods for dealing with Products
type Product struct {
	DB *sqlx.DB
}

// ListProducts tells you about request you made
func (p *Product) List(w http.ResponseWriter, r *http.Request) {
	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error querying DB : ", err)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling : ", err)
		return
	}

	w.Header().Set("content-type", "application/json; charset= utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(data); err != nil {
		log.Println("Error writing : ", err)
	}
}
