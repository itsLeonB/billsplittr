-- +goose Up
CREATE TABLE IF NOT EXISTS group_expense_bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payer_profile_id UUID NOT NULL REFERENCES user_profiles(id),
    image_name TEXT NOT NULL,
    group_expense_id UUID REFERENCES group_expenses(id),
    creator_profile_id UUID NOT NULL REFERENCES user_profiles(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS group_expense_bills;
