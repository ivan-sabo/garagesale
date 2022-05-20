package product_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ivan-sabo/garagesale/internal/platform/database/databasetest"
	"github.com/ivan-sabo/garagesale/internal/product"
)

func TestProducts(t *testing.T) {
	db, cleanup := databasetest.Setup(t)

	defer cleanup()

	np := product.NewProduct{
		Name:     "Comic Books",
		Cost:     10,
		Quantity: 20,
	}

	now := time.Now().UTC()

	p0, err := product.Create(db, np, now)
	if err != nil {
		t.Fatalf("could not create product: %v", err)
	}

	p1, err := product.Retrieve(db, p0.ID)
	if err != nil {
		t.Fatalf("could not retrive product: %v", err)
	}

	if diff := cmp.Diff(p0, p1); diff != "" {
		t.Fatalf("saved product did not match created: see diff \n%s", diff)
	}
}
