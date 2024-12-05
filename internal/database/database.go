package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Initialize(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		dbPath = "rockabilly.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("error enabling foreign keys: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("error creating tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS customers (
            id TEXT PRIMARY KEY,
            company_name TEXT NOT NULL,
            created_at DATETIME NOT NULL,
            status TEXT NOT NULL CHECK (status IN ('new', 'active', 'closed')),
            reference TEXT,
            internal_note TEXT,
            phone TEXT,
            tax_number TEXT,
            tax_rate_id TEXT,
            minimum_spend REAL,
            payment_terms_id TEXT,
            customer_group_id TEXT,
            price_list_id TEXT,
            order_interval INTEGER CHECK (order_interval IN (1,2,3,4)),
            email_addresses TEXT NOT NULL, -- JSON object
            buyers TEXT NOT NULL -- JSON array
        );`,

		`CREATE TABLE IF NOT EXISTS addresses (
            id TEXT PRIMARY KEY,
            customer_id TEXT NOT NULL,
            type TEXT NOT NULL CHECK (type IN ('shipping', 'billing')),
            company_name TEXT,
            contact_name TEXT,
            line1 TEXT NOT NULL,
            line2 TEXT,
            city TEXT NOT NULL,
            state TEXT,
            postal_code TEXT NOT NULL,
            country TEXT NOT NULL,
            FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
        );`,

		`CREATE TABLE IF NOT EXISTS orders (
            id TEXT PRIMARY KEY,
            number INTEGER NOT NULL UNIQUE,
            created DATETIME NOT NULL,
            status TEXT NOT NULL CHECK (status IN ('new', 'invoiced', 'released', 'part_fulfilled', 'preorder', 'fulfilled', 'standing_order', 'cancelled')),
            customer_id TEXT NOT NULL,
            company_name TEXT NOT NULL,
            phone TEXT,
            email_addresses TEXT NOT NULL, -- JSON object
            created_by TEXT NOT NULL,
            delivery_date DATETIME NOT NULL,
            reference TEXT,
            internal_note TEXT,
            customer_po_number TEXT,
            customer_note TEXT,
            standing_order_id TEXT,
            shipping_type TEXT,
            shipping_address_id TEXT,
            billing_address_id TEXT,
            currency TEXT NOT NULL,
            net_total REAL NOT NULL,
            gross_total REAL NOT NULL,
            FOREIGN KEY (customer_id) REFERENCES customers(id),
            FOREIGN KEY (shipping_address_id) REFERENCES addresses(id),
            FOREIGN KEY (billing_address_id) REFERENCES addresses(id)
        );`,

		`CREATE TABLE IF NOT EXISTS order_lines (
            id TEXT PRIMARY KEY,
            order_id TEXT NOT NULL,
            sku TEXT NOT NULL,
            name TEXT NOT NULL,
            options TEXT,
            grouping_category_id TEXT,
            grouping_category_name TEXT,
            shipping BOOLEAN NOT NULL DEFAULT 0,
            quantity INTEGER NOT NULL,
            unit_price REAL NOT NULL,
            sub_total REAL NOT NULL,
            tax_rate_id TEXT NOT NULL,
            tax_name TEXT NOT NULL,
            tax_rate REAL NOT NULL,
            tax_amount REAL NOT NULL,
            preorder_window_id TEXT,
            on_hold BOOLEAN NOT NULL DEFAULT 0,
            invoiced INTEGER NOT NULL DEFAULT 0,
            paid INTEGER NOT NULL DEFAULT 0,
            dispatched INTEGER NOT NULL DEFAULT 0,
            FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
        );`,

		`CREATE INDEX IF NOT EXISTS idx_customers_status ON customers(status);`,
		`CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);`,
		`CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);`,
		`CREATE INDEX IF NOT EXISTS idx_orders_delivery_date ON orders(delivery_date);`,
		`CREATE INDEX IF NOT EXISTS idx_order_lines_order_id ON order_lines(order_id);`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return fmt.Errorf("error creating table: %w", err)
		}
	}

	return nil
}
