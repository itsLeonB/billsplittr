-- +goose Up
ALTER TABLE group_expense_bills
ADD COLUMN extracted_text TEXT;

-- +goose Down
ALTER TABLE group_expense_bills
DROP COLUMN extracted_text;
