package product

import "github.com/jmoiron/sqlx"

// List returns all known products
func List(db *sqlx.DB) ([]Product, error) {
	list := []Product{}

	const q = `SELECT * FROM products`

	if err := db.Select(&list, q); err != nil {
		return nil, err
	}

	return list, nil
}
