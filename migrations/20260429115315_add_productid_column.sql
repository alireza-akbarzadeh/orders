-- +goose Up
ALTER TABLE orders
    ADD COLUMN product_id VARCHAR(36);

-- Optional: Add an index for faster lookups
CREATE INDEX idx_orders_product_id ON orders(product_id);

-- +goose Down
DROP INDEX IF EXISTS idx_orders_product_id ON orders;

ALTER TABLE orders
DROP COLUMN product_id;