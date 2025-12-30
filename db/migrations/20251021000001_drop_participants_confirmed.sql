-- +goose Up
ALTER TABLE IF EXISTS group_expenses
DROP COLUMN participants_confirmed;

ALTER TABLE IF EXISTS group_expense_participants
DROP COLUMN confirmed;

-- +goose Down
ALTER TABLE IF EXISTS group_expense_participants
ADD COLUMN confirmed BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE IF EXISTS group_expenses
ADD COLUMN participants_confirmed BOOLEAN NOT NULL DEFAULT FALSE;
