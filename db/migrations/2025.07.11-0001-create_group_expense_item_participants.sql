CREATE TABLE IF NOT EXISTS group_expense_item_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expense_item_id UUID NOT NULL REFERENCES group_expense_items(id),
    profile_id UUID NOT NULL REFERENCES user_profiles(id),
    share NUMERIC(20, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

ALTER TABLE group_expense_item_participants
ADD CONSTRAINT unique_expense_item_profile
UNIQUE (expense_item_id, profile_id);
