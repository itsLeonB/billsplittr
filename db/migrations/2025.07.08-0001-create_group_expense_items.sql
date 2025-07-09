CREATE TABLE IF NOT EXISTS group_expense_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_expense_id UUID NOT NULL REFERENCES group_expenses(id),
    name TEXT NOT NULL,
    amount NUMERIC(20, 2) NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS group_expense_other_fees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_expense_id UUID NOT NULL REFERENCES group_expenses(id),
    name TEXT NOT NULL,
    amount NUMERIC(20, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

ALTER TABLE group_expenses
ADD COLUMN confirmed BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN participants_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN created_by_profile_id UUID NOT NULL REFERENCES user_profiles(id);

ALTER TABLE group_expense_participants
ADD COLUMN confirmed BOOLEAN NOT NULL DEFAULT FALSE;
