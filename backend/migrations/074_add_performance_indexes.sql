-- Migration: Add performance indexes for slow queries
-- This migration adds indexes to improve query performance for frequently accessed data

-- Index for user_subscriptions: speeds up active subscription queries
-- Used by: ListActiveByUserID (user_id, status, expires_at)
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_user_status_expires
ON user_subscriptions(user_id, status, expires_at);

-- Index for user_subscriptions: speeds up group-based queries
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_group_status
ON user_subscriptions(group_id, status);

-- Index for usage_logs: speeds up user usage queries with date range
CREATE INDEX IF NOT EXISTS idx_usage_logs_user_created
ON usage_logs(user_id, created_at);

-- Index for usage_logs: speeds up API key usage queries
CREATE INDEX IF NOT EXISTS idx_usage_logs_api_key_created
ON usage_logs(api_key_id, created_at);

-- Index for api_keys: speeds up user's API key lookups
CREATE INDEX IF NOT EXISTS idx_api_keys_user_status
ON api_keys(user_id, status);

-- Composite index for user_subscriptions: covers most common query patterns
CREATE INDEX IF NOT EXISTS idx_user_subscriptions_composite
ON user_subscriptions(user_id, status, expires_at, group_id);
