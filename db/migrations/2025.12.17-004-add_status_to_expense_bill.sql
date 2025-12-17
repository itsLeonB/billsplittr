ALTER TABLE group_expense_bills
ADD COLUMN status TEXT NOT NULL DEFAULT 'PENDING';

UPDATE group_expense_bills
SET status = 'PENDING'
WHERE group_expense_id IS NULL;

UPDATE group_expense_bills
SET status = 'PARSED'
WHERE group_expense_id IS NOT NULL;
