package middleware

import (
	"expvar"
	"net/http"
	"runtime"

	"github.com/ivan-sabo/garagesale/internal/platform/web"
)

var m = struct {
	gr  *expvar.Int
	req *expvar.Int
	err *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
	err: expvar.NewInt("errors"),
}

// Metrics uprates program counters.
func Metrics() web.Middleware {

	// This is the actual middleware function to be executed
	f := func(before web.Handler) web.Handler {

		// wrap this handler around the next one provided.
		h := func(w http.ResponseWriter, r *http.Request) error {
			err := before(w, r)

			// Increment the request counter
			m.req.Add(1)

			// Update the count for the number of active goroutines every 100 requests
			if m.req.Value()%100 == 0 {
				m.gr.Set(int64(runtime.NumGoroutine()))
			}

			// return the error so it can be handled further up the chain
			return err
		}

		return h
	}

	return f
}
