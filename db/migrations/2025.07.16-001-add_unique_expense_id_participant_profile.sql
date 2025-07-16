ALTER TABLE group_expense_participants
ADD CONSTRAINT unique_expense_profile
UNIQUE (group_expense_id, participant_profile_id);
