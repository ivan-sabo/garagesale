package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ivan-sabo/garagesale/internal/platform/database/databasetest"
	"github.com/ivan-sabo/garagesale/internal/product"
	"github.com/ivan-sabo/garagesale/internal/schema"
)

func TestProducts(t *testing.T) {
	db, cleanup := databasetest.Setup(t)
	defer cleanup()

	ctx := context.Background()

	np := product.NewProduct{
		Name:     "Comic Books",
		Cost:     10,
		Quantity: 20,
	}

	now := time.Now().UTC()

	p0, err := product.Create(ctx, db, np, now)
	if err != nil {
		t.Fatalf("could not create product: %v", err)
	}

	p1, err := product.Retrieve(ctx, db, p0.ID)
	if err != nil {
		t.Fatalf("could not retrive product: %v", err)
	}

	if diff := cmp.Diff(p0, p1); diff != "" {
		t.Fatalf("saved product did not match created: see diff \n%s", diff)
	}
}

func TestProductList(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	ps, err := product.List(ctx, db)
	if err != nil {
		t.Fatalf("listing products: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected product list size %v, got %v", exp, got)
	}
}
