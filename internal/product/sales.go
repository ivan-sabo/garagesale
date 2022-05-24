package product

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AddSale records a sales transation for a single Product
func AddSale(ctx context.Context, db *sqlx.DB, ns NewSale, productID string, now time.Time) (*Sale, error) {
	s := Sale{
		ID:          uuid.New().String(),
		ProductID:   productID,
		Quantity:    ns.Quantity,
		Paid:        ns.Paid,
		DateCreated: now,
	}

	const q = `INSERT INTO sales
	(sale_id, product_id, quantity, paid, date_created)
	VALUES ($1, $2, $3, $4, $5)`

	_, err := db.ExecContext(ctx, q,
		s.ID, s.ProductID, s.Quantity,
		s.Paid, s.DateCreated,
	)
	if err != nil {
		return nil, fmt.Errorf("inserting sale: %w", err)
	}

	return &s, nil
}

// ListSales gives all Sales for a Product
func ListSales(ctx context.Context, db *sqlx.DB, productID string) ([]Sale, error) {
	sales := []Sale{}

	const q = `SELECT * FROM sales WHERE product_id = $1`
	if err := db.SelectContext(ctx, &sales, q, productID); err != nil {
		return nil, fmt.Errorf("selecting sales: %w", err)
	}

	return sales, nil
}
