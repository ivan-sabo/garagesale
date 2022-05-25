package middleware

import (
	"log"
	"net/http"

	"github.com/ivan-sabo/garagesale/internal/platform/web"
)

func Errors(log *log.Logger) web.Middleware {
	// This is the actual middleware function to be executed
	f := func(before web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {

			// Run the handler chain and catch any propagated error.
			if err := before(w, r); err != nil {

				// Log the error.
				log.Printf("Error : %v", err)

				// Respond to the error.
				if err := web.ResponError(w, err); err != nil {
					return err
				}
			}

			// Return nil to indicate the error has been handled
			return nil
		}

		return h
	}

	return f
}
