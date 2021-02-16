package schema

import (
	"github.com/jmoiron/sqlx"
)

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
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

// seeds is a string constant containing all of the queries needed to get the
// db seeded to a useful state for development.
//
// Note that database servers besides PostgreSQL may not support running
// multiple queries as part of the same execution so this single large constant
// may need to be broken up.
const seeds = `
-- Create admin and regular User with password "gophers"
INSERT INTO users (user_id, name, email, roles, password_hash, date_created, date_updated) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', 'admin@example.com', '{ADMIN,USER}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Gopher', 'user@example.com', '{USER}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;
	-- Create category"
	INSERT INTO categories (category_id, title, slug, image, parrent_id, description, meta_description, meta_title, meta_keywords, date_created, date_updated) VALUES
	('00000000-0000-0000-0000-000000000000', 'First category', 'first-category', 'link-to-image','00000000-0000-0000-0000-000000000000', '', '','','','2020-02-04 00:00:00', '2020-02-04 00:00:00')
	ON CONFLICT DO NOTHING;
	INSERT INTO brands
	(brand_id, title, slug, description, image, meta_description, meta_title, meta_keywords,  date_created, date_updated) VALUES
	('84fc7ad7-0f6c-4938-9cec-bb8f55953709', 'Brand Title', 'brand-title', 'description text', 'link-to-image', '','','', '2020-02-04 00:00:00', '2020-02-04 00:00:00')
	ON CONFLICT DO NOTHING;
	INSERT INTO products
	(product_id, title, slug, category_id, brand_id, price,  description, short_description, image, meta_description, meta_title, meta_keywords,  date_created, date_updated) VALUES
	('9097a8f9-c7c0-4e88-81da-72ec34a1dc79', 'Product Title', 'product-title', '00000000-0000-0000-0000-000000000000', '84fc7ad7-0f6c-4938-9cec-bb8f55953709', '3535.23', 'description text','', 'link-to-image', '','','','2020-02-04 00:00:00', '2020-02-04 00:00:00')
	ON CONFLICT DO NOTHING;
`

// DeleteAll runs the set of Drop-table queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func DeleteAll(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(deleteAll); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

// deleteAll is used to clean the database between tests.
const deleteAll = `
DELETE FROM users;
DELETE FROM categories;
DELETE FROM products;
DELETE FROM brands;
DELETE FROM articles;
DELETE FROM article_categories;
DELETE FROM slides;
`
