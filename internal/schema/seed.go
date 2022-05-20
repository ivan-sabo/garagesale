package schema

import "github.com/jmoiron/sqlx"

const seeds = `
INSERT INTO products (product_id, name, cost, quantity, date_created, date_updated) VALUES
('fb5c6c41-2b8a-499a-abd7-ab4d02bd2c01', 'Comic Books', 50, 42, '1999-01-08 04:05:06', '1999-01-08 04:05:06'),
('67621e3c-b845-4379-9ec8-875c8b2702c6', 'McDonalds Toys', 75, 120, '2020-04-04 04:05:06', '2020-04-04 04:05:06')
ON CONFLICT DO NOTHING;`

func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return tx.Commit()
}
