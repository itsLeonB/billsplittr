ALTER TABLE group_expense_item_participants
DROP CONSTRAINT group_expense_item_participants_expense_item_id_fkey,
ADD CONSTRAINT group_expense_item_participants_expense_item_id_fkey
  FOREIGN KEY (expense_item_id)
  REFERENCES group_expense_items(id)
  ON DELETE CASCADE;
