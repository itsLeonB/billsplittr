CREATE TYPE fee_calculation_method AS ENUM (
    'EQUAL_SPLIT',
    'ITEMIZED_SPLIT'
);

ALTER TABLE group_expense_other_fees
ADD COLUMN IF NOT EXISTS calculation_method fee_calculation_method;

CREATE TABLE IF NOT EXISTS group_expense_other_fee_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    other_fee_id UUID NOT NULL REFERENCES group_expense_other_fees(id),
    profile_id UUID NOT NULL REFERENCES user_profiles(id),
    share_amount NUMERIC(20, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

ALTER TABLE group_expense_other_fees
ADD COLUMN IF NOT EXISTS rate NUMERIC(20, 4);
