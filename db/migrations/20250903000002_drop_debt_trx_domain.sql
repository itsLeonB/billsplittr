-- +goose Up
DROP TABLE IF EXISTS debt_transactions;
DROP TABLE IF EXISTS transfer_methods;

-- +goose Down
-- This migration is destructive and cannot be easily reversed.
