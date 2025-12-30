-- +goose Up
ALTER TABLE IF EXISTS group_expenses
ADD COLUMN IF NOT EXISTS subtotal NUMERIC(20, 2);

-- +goose Down
ALTER TABLE IF EXISTS group_expenses
DROP COLUMN IF EXISTS subtotal;
