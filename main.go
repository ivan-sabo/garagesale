package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Printf("Main : started")
	defer log.Println("Main : completed")

	api := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(Echo),
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

// Echo tells you about request you made
func Echo(w http.ResponseWriter, r *http.Request) {
	id := rand.Intn(1000)

	fmt.Println("starting", id)

	time.Sleep(3 * time.Second)

	fmt.Fprintln(w, "You asked to ", r.Method, r.URL.Path)

	fmt.Println("Ending", id)
}
