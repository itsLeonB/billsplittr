-- +goose Up
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_profiles;
DROP TABLE IF EXISTS friendships;

-- +goose Down
-- This migration is destructive and cannot be easily reversed.
