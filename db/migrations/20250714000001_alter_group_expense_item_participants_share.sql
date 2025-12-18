-- +goose Up
ALTER TABLE group_expense_item_participants
ALTER COLUMN share TYPE NUMERIC(20, 4);

-- +goose Down
ALTER TABLE group_expense_item_participants
ALTER COLUMN share TYPE NUMERIC(20, 2);
