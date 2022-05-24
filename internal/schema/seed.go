package schema

import "github.com/jmoiron/sqlx"

const seeds = `
INSERT INTO products (product_id, name, cost, quantity, date_created, date_updated) VALUES
('fb5c6c41-2b8a-499a-abd7-ab4d02bd2c01', 'Comic Books', 50, 42, '1999-01-08 04:05:06', '1999-01-08 04:05:06'),
('67621e3c-b845-4379-9ec8-875c8b2702c6', 'McDonalds Toys', 75, 120, '2020-04-04 04:05:06', '2020-04-04 04:05:06')
ON CONFLICT DO NOTHING;

INSERT INTO sales (sale_id, product_id, quantity, paid, date_created) VALUES
	('dc3ea3fa-dcfc-4073-8fa1-7187d44eaa14', 'fb5c6c41-2b8a-499a-abd7-ab4d02bd2c01', 2, 100, '2021-01-18 14:05:06'),
	('bf27a541-e746-4762-a3dc-641f86e3e06c', 'fb5c6c41-2b8a-499a-abd7-ab4d02bd2c01', 4, 300, '2015-06-12 06:05:06')
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
