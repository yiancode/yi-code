-- Migration: Add user usage report settings
-- Description: Adds fields for user-level usage report email configuration

-- Add usage report configuration fields to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS usage_report_enabled BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS usage_report_schedule VARCHAR(20) NOT NULL DEFAULT '09:00',
ADD COLUMN IF NOT EXISTS usage_report_timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Shanghai';

-- Create index for efficient querying of users with reports enabled
CREATE INDEX IF NOT EXISTS idx_users_usage_report_enabled ON users(usage_report_enabled) WHERE deleted_at IS NULL;
