package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // Register the /debug/pprof handlers
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/ivan-sabo/garagesale/cmd/sales-api/internal/handlers"
	"github.com/ivan-sabo/garagesale/internal/platform/database"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Shadowing - overriding log package with local variable
	log := log.New(os.Stdout, "sales :", log.LstdFlags)

	var cfg struct {
		Web struct {
			Address         string        `env:"ADDRESS" envDefault:"localhost:8000"`
			Debug           string        `env:"DEBUG" envDefault:"localhost:6060"`
			ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
			WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" envDefault:"5s"`
			ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
		}
		DB struct {
			User       string `env:"USER" envDefault:"postgres"`
			Password   string `env:"PASSWORD" envDefault:"postgres,noprint"`
			Host       string `env:"HOST" envDefault:"localhost"`
			Name       string `env:"NAME" envDefault:"postgres"`
			DisableTLS bool   `env:"DISABLE_TLS" envDefault:"true"`
		}
	}

	log.Printf("Main : started")
	defer log.Println("Main : completed")

	// GetConfiguration
	if err := env.Parse(&cfg, env.Options{Prefix: "SALE_"}); err != nil {
		log.Fatalf("error: parsing config: %s", err)
	}

	// Setup dependencies
	db, err := database.Open(database.Config{
		Host:       cfg.DB.Host,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
		Name:       cfg.DB.Name,
	})
	if err != nil {
		return fmt.Errorf("run database error : %w", err)
	}
	defer db.Close()

	// Start Debug service
	go func() {
		log.Printf("main : Debug service listening on %s", cfg.Web.Debug)
		err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
		log.Printf("main: Debug service ended %v", err)
	}()

	// Start API service
	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      handlers.API(log, db),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
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
		return fmt.Errorf("error : listening and serving : %w", err)
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
				return fmt.Errorf("main : could not stop server gracefully: %w", err)
			}
		}
	}

	return nil
}
