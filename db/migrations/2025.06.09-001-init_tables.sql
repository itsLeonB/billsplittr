-- Enums
CREATE TYPE debt_transaction_type AS ENUM ('LEND', 'REPAY');
CREATE TYPE friendship_type AS ENUM ('REAL', 'ANON');

-- Tables
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
COMMENT ON COLUMN user_profiles.user_id IS 'Nullable. Can be NULL for peers who do not have an account in the app';

CREATE TABLE IF NOT EXISTS friendships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id_1 UUID NOT NULL REFERENCES user_profiles(id),
    profile_id_2 UUID NOT NULL REFERENCES user_profiles(id),
    type friendship_type NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT unique_friendship UNIQUE (profile_id_1, profile_id_2),
    CONSTRAINT profile_order CHECK (profile_id_1 < profile_id_2)
);

CREATE TABLE IF NOT EXISTS transfer_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    display TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS debt_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lender_profile_id UUID NOT NULL REFERENCES user_profiles(id),
    borrower_profile_id UUID NOT NULL REFERENCES user_profiles(id),
    type debt_transaction_type NOT NULL,
    amount NUMERIC(20, 2) NOT NULL,
    transfer_method_id UUID NOT NULL REFERENCES transfer_methods(id),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS group_expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payer_profile_id UUID NOT NULL REFERENCES user_profiles(id),
    total_amount NUMERIC(20, 2) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS group_expense_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_expense_id UUID NOT NULL REFERENCES group_expenses(id),
    participant_profile_id UUID NOT NULL REFERENCES user_profiles(id),
    share_amount NUMERIC(20, 2) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Indexes
CREATE INDEX IF NOT EXISTS user_profiles_user_id_idx ON user_profiles(user_id);
CREATE INDEX IF NOT EXISTS user_profiles_name_idx ON user_profiles(name);
CREATE INDEX IF NOT EXISTS friendships_profile_id_1_idx ON friendships(profile_id_1);
CREATE INDEX IF NOT EXISTS friendships_profile_id_2_idx ON friendships(profile_id_2);
CREATE INDEX IF NOT EXISTS friendships_type_idx ON friendships(type);
CREATE INDEX IF NOT EXISTS debt_transactions_lender_profile_id_idx ON debt_transactions(lender_profile_id);
CREATE INDEX IF NOT EXISTS debt_transactions_borrower_profile_id_idx ON debt_transactions(borrower_profile_id);
CREATE INDEX IF NOT EXISTS debt_transactions_transfer_method_id_idx ON debt_transactions(transfer_method_id);
CREATE INDEX IF NOT EXISTS debt_transactions_created_at_idx ON debt_transactions(created_at);
CREATE INDEX IF NOT EXISTS group_expenses_payer_profile_id_idx ON group_expenses(payer_profile_id);
CREATE INDEX IF NOT EXISTS group_expenses_created_at_idx ON group_expenses(created_at);
CREATE INDEX IF NOT EXISTS group_expense_participants_group_expense_id_idx ON group_expense_participants(group_expense_id);
CREATE INDEX IF NOT EXISTS group_expense_participants_participant_profile_id_idx ON group_expense_participants(participant_profile_id);
CREATE INDEX IF NOT EXISTS group_expense_participants_created_at_idx ON group_expense_participants(created_at);
