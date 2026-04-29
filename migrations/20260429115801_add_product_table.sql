-- +goose Up
CREATE TABLE IF NOT EXISTS products (
                                        id TEXT PRIMARY KEY,
                                        name TEXT NOT NULL,
                                        description TEXT,
                                        price DECIMAL(12,2) NOT NULL,
    stock_quantity INTEGER NOT NULL DEFAULT 0,
    sku TEXT UNIQUE,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
    );

-- +goose Down
DROP TABLE IF EXISTS products;
