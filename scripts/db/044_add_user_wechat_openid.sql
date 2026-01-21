-- Migration: Add wechat_openid field to users table for WeChat account binding
-- This allows users to bind their WeChat account for QR code login

-- +goose Up
-- +goose StatementBegin

-- Add wechat_openid column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS wechat_openid VARCHAR(64) DEFAULT '';

-- Create unique partial index (only for non-empty, non-deleted values)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_wechat_openid_unique
    ON users(wechat_openid) WHERE wechat_openid != '' AND deleted_at IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_users_wechat_openid_unique;
ALTER TABLE users DROP COLUMN IF EXISTS wechat_openid;

-- +goose StatementEnd
