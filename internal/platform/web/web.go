package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// App is the entrypoint for all web applications
type App struct {
	mux *chi.Mux
	log *log.Logger
}

// NewApp knows how to construct internal state for an App
func NewApp(l *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(),
		log: l,
	}
}

// Handle connects a method and URL pattern to a particular application handler
func (a *App) Handle(method, pattern string, fn http.HandlerFunc) {
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}