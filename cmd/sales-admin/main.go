package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/ivan-sabo/garagesale/internal/platform/database"
	"github.com/ivan-sabo/garagesale/internal/schema"
	"github.com/jmoiron/sqlx"
)

func main() {
	// Setup dependencies
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("applying migrations: ", err)
		}
		log.Println("Migrations complete")
		return
	}
}

// Product is something we sell
type Product struct {
	ID          string    `db:"product_id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Cost        int       `db:"cost" json:"cost"`
	Quantity    int       `db:"quantity" json:"quantity"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// ProductService has handler methods for dealing with Products
type ProductService struct {
	db *sqlx.DB
}

// ListProducts tells you about request you made
func (p *ProductService) List(w http.ResponseWriter, r *http.Request) {
	list := []Product{}

	const q = `SELECT * FROM products`

	if err := p.db.Select(&list, q); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error querying db", err)
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
