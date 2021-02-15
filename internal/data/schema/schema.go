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
	image TEXT,
	description TEXT,
	meta_title          TEXT,
	meta_keywords          TEXT,
	meta_description          TEXT,
	date_created  TIMESTAMP,
	date_updated  TIMESTAMP,

	PRIMARY KEY (category_id),
	FOREIGN KEY (parrent_id) REFERENCES categories(category_id) ON DELETE SET NULL

	);`,
	},
	{
		Version:     1.2,
		Description: "Create table Brands",
		Script: `
CREATE TABLE brands (
	brand_id       UUID,
	title          TEXT,
	slug         TEXT UNIQUE,
	image TEXT,
	description TEXT,
	meta_title          TEXT,
	meta_keywords          TEXT,
	meta_description          TEXT,
	date_created  TIMESTAMP,
	date_updated  TIMESTAMP,

	PRIMARY KEY (brand_id)

	);`,
	},
	{
		Version:     1.3,
		Description: "Create table Products",
		Script: `
		CREATE TABLE products (
			product_id       UUID,
			title          TEXT NOT NULL,
			slug         TEXT UNIQUE NOT NULL,
			category_id   UUID,
			brand_id   UUID,
			price NUMERIC(15,2) NOT NULL DEFAULT 0.00,
			old_price NUMERIC(15,2) NOT NULL DEFAULT 0.00,
			short_description TEXT,
			description TEXT,
			image TEXT,
			meta_title          TEXT,
			meta_keywords          TEXT,
			meta_description          TEXT,
			date_created  TIMESTAMP,
			date_updated  TIMESTAMP,
		
			PRIMARY KEY (product_id),
			FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE SET NULL,
			FOREIGN KEY (brand_id) REFERENCES brands(brand_id) ON DELETE SET NULL
			);`,
	},
	{
		Version:     1.4,
		Description: "Create table Slides",
		Script: `
CREATE TABLE Slides (
	slide_id       UUID,
	title          TEXT,
	sub_title TEXT,
	image TEXT,
	link TEXT,	
	date_created  TIMESTAMP,
	date_updated  TIMESTAMP,

	PRIMARY KEY (slide_id)
	);`,
	},
	{
		Version:     1.5,
		Description: "Create table Article Categories",
		Script: `
CREATE TABLE article_categories (
	category_id       UUID,
	title          TEXT,
	slug         TEXT UNIQUE,
	image TEXT,
	description TEXT,
	meta_title          TEXT,
	meta_keywords          TEXT,
	meta_description          TEXT,
	date_created  TIMESTAMP,
	date_updated  TIMESTAMP,

	PRIMARY KEY (category_id)
	);`,
	},
	{
		Version:     1.6,
		Description: "Create table Articles",
		Script: `
CREATE TABLE articles (
	article_id       UUID,
	title          TEXT,
	slug         TEXT UNIQUE,
	image TEXT,
	category_id   UUID,
	description TEXT,
	meta_title          TEXT,
	meta_keywords          TEXT,
	meta_description          TEXT,
	date_created  TIMESTAMP,
	date_updated  TIMESTAMP,

	PRIMARY KEY (article_id),
	FOREIGN KEY (category_id) REFERENCES article_categories(category_id) ON DELETE SET NULL
	);`,
	},
}
