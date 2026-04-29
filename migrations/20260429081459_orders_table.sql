-- +goose Up
-- +goose statementbegin
ALTER TABLE orders ADD COLUMN tracking_number TEXT;
ALTER TABLE orders ADD COLUMN cancelled_at DATETIME;
CREATE UNIQUE INDEX idx_tracking_number ON orders(tracking_number);
-- +goose statementend

-- +goose Down
-- +goose statementbegin
DROP INDEX IF EXISTS idx_tracking_number;
ALTER TABLE orders DROP COLUMN cancelled_at;
ALTER TABLE orders DROP COLUMN tracking_number;
-- +goose statementend