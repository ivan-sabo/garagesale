package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	log.Printf("Main : started")
	defer log.Println("Main : completed")

	// Setup dependencies
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	api := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(ListProducts),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("main : API listening on %s", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("error : listening and serving: %s", err)
	case <-shutdown:
		log.Println("main : Start shutdown")

		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v: %v", timeout, err)

			err := api.Close()
			if err != nil {
				log.Fatalf("main : could not stop server gracefully: %v", err)
			}
		}
	}

}

func openDB() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmofe", "disabled")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

// Product is something we sell
type Product struct {
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Quantity int    `json:"quantity"`
}

// ListProducts tells you about request you made
func ListProducts(w http.ResponseWriter, r *http.Request) {
	list := []Product{
		{Name: "Comic books", Cost: 75, Quantity: 50},
		{Name: "McDonalds toys", Cost: 25, Quantity: 125},
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
