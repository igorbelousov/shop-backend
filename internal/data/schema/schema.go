// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

var migrations = []darwin.Migration{
	{
		Version:     1.0,
		Description: "Create table users",
		Script: `
CREATE TABLE users (
	user_id       UUID,
	name          TEXT,
	email         TEXT UNIQUE,
	roles         TEXT[],
	password_hash TEXT,
	date_created  TIMESTAMP,
	date_updated  TIMESTAMP,

	PRIMARY KEY (user_id)
);`,
	},
	{
		Version:     1.1,
		Description: "Create table Categories",
		Script: `
CREATE TABLE categories (
	category_id       UUID,
	title          TEXT,
	slug         TEXT UNIQUE,
	parrent_id   UUID,
	description TEXT,
	date_created  TIMESTAMP,
	date_updated  TIMESTAMP,

	PRIMARY KEY (category_id),
	FOREIGN KEY (parrent_id) REFERENCES categories(category_id) ON DELETE SET NULL

);`,
	},
	{
		Version:     1.2,
		Description: "Create table Products",
		Script: `
		CREATE TABLE products (
			product_id       UUID,
			title          TEXT NOT NULL,
			slug         TEXT UNIQUE NOT NULL,
			category_id   UUID,
			price NUMERIC(15,2) NOT NULL DEFAULT 0.00,
			description TEXT,
			date_created  TIMESTAMP,
			date_updated  TIMESTAMP,
		
			PRIMARY KEY (product_id),
			FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE SET NULL
		
		);`,
	},
}
