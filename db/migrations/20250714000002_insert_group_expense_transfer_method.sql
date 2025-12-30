-- +goose Up
INSERT INTO transfer_methods (name, display)
VALUES ('GROUP_EXPENSE', 'Group Expense');

-- +goose Down
DELETE FROM transfer_methods WHERE name = 'GROUP_EXPENSE';
