ALTER TABLE group_expenses
ADD COLUMN status TEXT NOT NULL DEFAULT 'DRAFT',
ADD COLUMN items_total NUMERIC(20, 2) NOT NULL DEFAULT 0,
ADD COLUMN fees_total NUMERIC(20, 2) NOT NULL DEFAULT 0;

UPDATE group_expenses
SET status = 'DRAFT'
WHERE confirmed IS FALSE;

UPDATE group_expenses
SET status = 'CONFIRMED'
WHERE confirmed IS TRUE;
