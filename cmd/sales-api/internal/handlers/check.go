package handlers

import (
	"net/http"

	"github.com/ivan-sabo/garagesale/internal/platform/database"
	"github.com/ivan-sabo/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// Check has handlers to implement service orchestration.
type Check struct {
	db *sqlx.DB
}

// Health respons with a 200 OK if the service is healthy and ready for traffic.
func (c *Check) Health(w http.ResponseWriter, r *http.Request) error {

	var health struct {
		Status string `json:"status"`
	}
	if err := database.StatusCheck(r.Context(), c.db); err != nil {
		health.Status = "db not ready"
		return web.Respond(w, health, http.StatusInternalServerError)
	}

	health.Status = "OK"
	return web.Respond(w, health, http.StatusOK)
}
