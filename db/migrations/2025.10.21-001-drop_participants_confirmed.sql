ALTER TABLE IF EXISTS group_expenses
DROP COLUMN participants_confirmed;

ALTER TABLE IF EXISTS group_expense_participants
DROP COLUMN confirmed;
