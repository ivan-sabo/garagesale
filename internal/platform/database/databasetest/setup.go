package databasetest

import (
	"testing"
	"time"

	"github.com/ivan-sabo/garagesale/internal/platform/database"
	"github.com/ivan-sabo/garagesale/internal/schema"
	"github.com/jmoiron/sqlx"
)

func Setup(t *testing.T) (*sqlx.DB, func()) {
	// should be called in every helper function - tells which line of code failed in code
	t.Helper()

	c := startContainer(t)

	db, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("opening database connection: %v", err)
	}

	t.Log("waiting for database to be ready")

	// Wait for the database to be ready. Wait 100ms longer between each attempt.
	// Do not try more than 20 times
	var pingError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		stopContainer(t, c)
		t.Fatalf("waiting for database to be ready: %v", pingError)
	}

	if err := schema.Migrate(db); err != nil {
		stopContainer(t, c)
		t.Fatalf("migrating: %s", err)
	}

	// teardown is the function that should be invoked when the caller is done
	//with the database
	teardown := func() {
		t.Helper()
		db.Close()
		stopContainer(t, c)
	}

	return db, teardown
}
