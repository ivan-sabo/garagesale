package database

import (
	"context"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Register the postgres database/sql driver
)

// Config is what we require to open a database connection
type Config struct {
	Host       string
	Name       string
	User       string
	Password   string
	DisableTLS bool
}

// Open knows how to open a database connection
func Open(cfg Config) (*sqlx.DB, error) {
	q := url.Values{}

	q.Set("sslmode", "require")
	if cfg.DisableTLS {
		q.Set("sslmode", "disable")
	}
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

// StatusCheck returns nil if it can successfully talk to the database.
// It returns a non-nil error otherwise
func StatusCheck(ctx context.Context, db *sqlx.DB) error {

	// return db.Ping()
	// will cache a connections. If DB goes away after it was previosly able to make a connection, this wouldn't be correct anymore.
	// SQL-specific thing. The only way to make sure is to run a real query
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
