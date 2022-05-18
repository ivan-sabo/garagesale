package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ivan-sabo/garagesale/cmd/sales-api/internal/handlers"
	"github.com/ivan-sabo/garagesale/internal/platform/database"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Printf("Main : started")
	defer log.Println("Main : completed")

	// Setup dependencies

	db, err := database.Open(database.Config{
		Host:       "localhost",
		User:       "postgres",
		Password:   "postgres",
		DisableTLS: true,
		Name:       "postgres",
	})
	if err != nil {
		return fmt.Errorf("Run database error : %w", err)
	}
	defer db.Close()

	ps := handlers.Product{DB: db}

	api := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(ps.List),
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
		fmt.Errorf("error : listening and serving : %w", err)
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
				fmt.Errorf("main : could not stop server gracefully: %w", err)
			}
		}
	}

	return nil
}
